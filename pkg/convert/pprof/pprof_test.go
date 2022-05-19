package pprof

import (
	"context"
	"sort"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/pyroscope-io/pyroscope/pkg/storage/tree"
	"github.com/pyroscope-io/pyroscope/pkg/storage/types"
)

type mockIngester struct{ actual []*types.PutInput }

func (m *mockIngester) Enqueue(_ context.Context, p *types.PutInput) {
	m.actual = append(m.actual, p)
}

var _ = Describe("pprof parsing", func() {
	Context("Go", func() {
		It("can parse CPU profiles", func() {
			p, err := readPprofFixture("testdata/cpu.pb.gz")
			Expect(err).ToNot(HaveOccurred())

			ingester := new(mockIngester)
			spyName := "spy-name"
			now := time.Now()
			start := now
			end := now.Add(10 * time.Second)
			labels := map[string]string{
				"__name__": "app",
				"foo":      "bar",
			}

			w := NewProfileWriter(ingester, ProfileWriterConfig{
				SampleTypes: tree.DefaultSampleTypeMapping,
				Labels:      labels,
				SpyName:     spyName,
			})

			err = w.WriteProfile(context.Background(), start, end, p)
			Expect(err).ToNot(HaveOccurred())

			Expect(ingester.actual).To(HaveLen(1))
			input := ingester.actual[0]
			Expect(input.SpyName).To(Equal(spyName))
			Expect(input.StartTime).To(Equal(start))
			Expect(input.EndTime).To(Equal(end))
			Expect(input.SampleRate).To(Equal(uint32(100)))
			Expect(input.Val.Samples()).To(Equal(uint64(47)))
			Expect(input.Key.Normalized()).To(Equal("app.cpu{foo=bar}"))
			Expect(input.Val.String()).To(ContainSubstring("runtime.main;main.main;main.slowFunction;main.work 1"))
		})
	})

	Context("JS", func() {
		It("can parses CPU profile", func() {
			p, err := readPprofFixture("testdata/nodejs-wall.pb.gz")
			Expect(err).ToNot(HaveOccurred())

			ingester := new(mockIngester)
			spyName := "nodespy"
			now := time.Now()
			start := now
			end := now.Add(10 * time.Second)
			labels := map[string]string{
				"__name__": "app",
				"foo":      "bar",
			}

			w := NewProfileWriter(ingester, ProfileWriterConfig{
				SampleTypes: tree.DefaultSampleTypeMapping,
				Labels:      labels,
				SpyName:     spyName,
			})

			err = w.WriteProfile(context.Background(), start, end, p)
			Expect(err).ToNot(HaveOccurred())

			Expect(ingester.actual).To(HaveLen(1))
			input := ingester.actual[0]
			Expect(input.SpyName).To(Equal(spyName))
			Expect(input.StartTime).To(Equal(start))
			Expect(input.EndTime).To(Equal(end))
			Expect(input.SampleRate).To(Equal(uint32(100)))
			Expect(input.Val.Samples()).To(Equal(uint64(898)))
			Expect(input.Key.Normalized()).To(Equal("app.cpu{foo=bar}"))
			Expect(input.Val.String()).To(ContainSubstring("node:_http_server:resOnFinish:819;node:_http_server:detachSocket:252 1"))
		})

		It("can parse heap profiles", func() {
			p, err := readPprofFixture("testdata/nodejs-heap.pb.gz")
			Expect(err).ToNot(HaveOccurred())

			ingester := new(mockIngester)
			spyName := "nodespy"
			now := time.Now()
			start := now
			end := now.Add(10 * time.Second)
			labels := map[string]string{
				"__name__": "app",
				"foo":      "bar",
			}

			Expect(tree.DefaultSampleTypeMapping["inuse_objects"].Cumulative).To(BeFalse())
			Expect(tree.DefaultSampleTypeMapping["inuse_space"].Cumulative).To(BeFalse())
			tree.DefaultSampleTypeMapping["inuse_objects"].Cumulative = false
			tree.DefaultSampleTypeMapping["inuse_space"].Cumulative = false

			w := NewProfileWriter(ingester, ProfileWriterConfig{
				SampleTypes: tree.DefaultSampleTypeMapping,
				Labels:      labels,
				SpyName:     spyName,
			})

			err = w.WriteProfile(context.Background(), start, end, p)
			Expect(err).ToNot(HaveOccurred())
			Expect(ingester.actual).To(HaveLen(2))
			sort.Slice(ingester.actual, func(i, j int) bool {
				return ingester.actual[i].Key.Normalized() < ingester.actual[j].Key.Normalized()
			})

			input := ingester.actual[0]
			Expect(input.SpyName).To(Equal(spyName))
			Expect(input.StartTime).To(Equal(start))
			Expect(input.EndTime).To(Equal(end))
			Expect(input.Val.Samples()).To(Equal(uint64(100498)))
			Expect(input.Key.Normalized()).To(Equal("app.inuse_objects{foo=bar}"))
			Expect(input.Val.String()).To(ContainSubstring("node:internal/streams/readable:readableAddChunk:236 138"))

			input = ingester.actual[1]
			Expect(input.SpyName).To(Equal(spyName))
			Expect(input.StartTime).To(Equal(start))
			Expect(input.EndTime).To(Equal(end))
			Expect(input.Val.Samples()).To(Equal(uint64(8357762)))
			Expect(input.Key.Normalized()).To(Equal("app.inuse_space{foo=bar}"))
			Expect(input.Val.String()).To(ContainSubstring("node:internal/net:isIPv6:35;:test:0 555360"))
		})
	})
})

var _ = Describe("pprof profile_id multiplexing", func() {
	It("can parse profiles labeled with profile_id correctly", func() {
		p, err := readPprofFixture("testdata/cpu-mux.pb.gz")
		Expect(err).ToNot(HaveOccurred())

		ingester := new(mockIngester)
		spyName := "spy-name"
		now := time.Now()
		start := now
		end := now.Add(10 * time.Second)
		labels := map[string]string{
			"__name__": "app",
			"foo":      "bar",
		}

		w := NewProfileWriter(ingester, ProfileWriterConfig{
			SampleTypes: tree.DefaultSampleTypeMapping,
			Labels:      labels,
			SpyName:     spyName,
		})

		err = w.WriteProfile(context.Background(), start, end, p)
		Expect(err).ToNot(HaveOccurred())

		var actualTotal uint64
		const (
			expectedTotal = uint64(789)
			expectedDiff  = uint64(20)
		)

		for _, v := range ingester.actual {
			if v.Key.Normalized() == "app.cpu{foo=bar}" {
				Expect(v.Val.Samples()).To(Equal(expectedTotal))
				continue
			}
			actualTotal += v.Val.Samples()
		}

		Expect(expectedTotal - actualTotal).To(Equal(expectedDiff))
	})
})

var _ = Describe("custom pprof parsing", func() {
	It("parses data correctly", func() {
		p, err := readPprofFixture("testdata/heap-js.pprof")
		Expect(err).ToNot(HaveOccurred())

		ingester := new(mockIngester)
		spyName := "spy-name"
		now := time.Now()
		start := now
		end := now.Add(10 * time.Second)
		labels := map[string]string{
			"__name__": "app",
			"foo":      "bar",
		}

		w := NewProfileWriter(ingester, ProfileWriterConfig{
			SampleTypes: map[string]*tree.SampleTypeConfig{
				"objects": {
					Units:       "objects",
					Aggregation: "average",
				},
				"space": {
					Units:       "bytes",
					Aggregation: "average",
				},
			},
			Labels:  labels,
			SpyName: spyName,
		})

		err = w.WriteProfile(context.TODO(), start, end, p)
		Expect(err).ToNot(HaveOccurred())
		Expect(ingester.actual).To(HaveLen(2))
		sort.Slice(ingester.actual, func(i, j int) bool {
			return ingester.actual[i].Key.Normalized() < ingester.actual[j].Key.Normalized()
		})

		input := ingester.actual[0]
		Expect(input.SpyName).To(Equal(spyName))
		Expect(input.StartTime).To(Equal(start))
		Expect(input.EndTime).To(Equal(end))
		Expect(input.Val.Samples()).To(Equal(uint64(66148)))
		Expect(input.Key.Normalized()).To(Equal("app.objects{foo=bar}"))
		Expect(input.Val.String()).To(ContainSubstring("parserOnHeadersComplete;parserOnIncoming 2428"))

		input = ingester.actual[1]
		Expect(input.SpyName).To(Equal(spyName))
		Expect(input.StartTime).To(Equal(start))
		Expect(input.EndTime).To(Equal(end))
		Expect(input.Val.Samples()).To(Equal(uint64(6388384)))
		Expect(input.Key.Normalized()).To(Equal("app.space{foo=bar}"))
		Expect(input.Val.String()).To(ContainSubstring("parserOnHeadersComplete;parserOnIncoming 524448"))
	})
})
