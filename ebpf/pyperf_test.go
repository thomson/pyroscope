package ebpfspy

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPythonProcInfo(t *testing.T) {
	var maps string
	maps = `55bb7ceb9000-55bb7ceba000 r--p 00000000 fd:01 57017250                   /home/korniltsev/.asdf/installs/python/3.6.0/bin/python3.6
55bb7ceba000-55bb7cebb000 r-xp 00001000 fd:01 57017250                   /home/korniltsev/.asdf/installs/python/3.6.0/bin/python3.6
55bb7cebb000-55bb7cebc000 r--p 00002000 fd:01 57017250                   /home/korniltsev/.asdf/installs/python/3.6.0/bin/python3.6
55bb7cebc000-55bb7cebd000 r--p 00002000 fd:01 57017250                   /home/korniltsev/.asdf/installs/python/3.6.0/bin/python3.6
55bb7cebd000-55bb7cebe000 rw-p 00003000 fd:01 57017250                   /home/korniltsev/.asdf/installs/python/3.6.0/bin/python3.6
55bb7df87000-55bb7e02a000 rw-p 00000000 00:00 0                          [heap]
7f94a5c00000-7f94a6180000 r--p 00000000 fd:01 45623802                   /usr/lib/locale/locale-archive
7f94a619b000-7f94a6200000 rw-p 00000000 00:00 0 
7f94a6200000-7f94a6222000 r--p 00000000 fd:01 45643280                   /usr/lib/x86_64-linux-gnu/libc.so.6
7f94a6222000-7f94a639a000 r-xp 00022000 fd:01 45643280                   /usr/lib/x86_64-linux-gnu/libc.so.6
7f94a639a000-7f94a63f2000 r--p 0019a000 fd:01 45643280                   /usr/lib/x86_64-linux-gnu/libc.so.6
7f94a63f2000-7f94a63f6000 r--p 001f1000 fd:01 45643280                   /usr/lib/x86_64-linux-gnu/libc.so.6
7f94a63f6000-7f94a63f8000 rw-p 001f5000 fd:01 45643280                   /usr/lib/x86_64-linux-gnu/libc.so.6
7f94a63f8000-7f94a6405000 rw-p 00000000 00:00 0 
7f94a6417000-7f94a6517000 rw-p 00000000 00:00 0 
7f94a6517000-7f94a6525000 r--p 00000000 fd:01 45643925                   /usr/lib/x86_64-linux-gnu/libm.so.6
7f94a6525000-7f94a65a3000 r-xp 0000e000 fd:01 45643925                   /usr/lib/x86_64-linux-gnu/libm.so.6
7f94a65a3000-7f94a65fe000 r--p 0008c000 fd:01 45643925                   /usr/lib/x86_64-linux-gnu/libm.so.6
7f94a65fe000-7f94a65ff000 r--p 000e6000 fd:01 45643925                   /usr/lib/x86_64-linux-gnu/libm.so.6
7f94a65ff000-7f94a6600000 rw-p 000e7000 fd:01 45643925                   /usr/lib/x86_64-linux-gnu/libm.so.6
7f94a6600000-7f94a665c000 r--p 00000000 fd:01 57017251                   /home/korniltsev/.asdf/installs/python/3.6.0/lib/libpython3.6m.so.1.0
7f94a665c000-7f94a67f3000 r-xp 0005c000 fd:01 57017251                   /home/korniltsev/.asdf/installs/python/3.6.0/lib/libpython3.6m.so.1.0
7f94a67f3000-7f94a6886000 r--p 001f3000 fd:01 57017251                   /home/korniltsev/.asdf/installs/python/3.6.0/lib/libpython3.6m.so.1.0
7f94a6886000-7f94a6889000 r--p 00285000 fd:01 57017251                   /home/korniltsev/.asdf/installs/python/3.6.0/lib/libpython3.6m.so.1.0
7f94a6889000-7f94a68ef000 rw-p 00288000 fd:01 57017251                   /home/korniltsev/.asdf/installs/python/3.6.0/lib/libpython3.6m.so.1.0
7f94a68ef000-7f94a6920000 rw-p 00000000 00:00 0 
7f94a693e000-7f94a69c0000 rw-p 00000000 00:00 0 
7f94a69d3000-7f94a69da000 r--s 00000000 fd:01 45645140                   /usr/lib/x86_64-linux-gnu/gconv/gconv-modules.cache
7f94a69da000-7f94a69dc000 rw-p 00000000 00:00 0 
7f94a69dc000-7f94a69dd000 r--p 00000000 fd:01 45642507                   /usr/lib/x86_64-linux-gnu/ld-linux-x86-64.so.2
7f94a69dd000-7f94a6a05000 r-xp 00001000 fd:01 45642507                   /usr/lib/x86_64-linux-gnu/ld-linux-x86-64.so.2
7f94a6a05000-7f94a6a0f000 r--p 00029000 fd:01 45642507                   /usr/lib/x86_64-linux-gnu/ld-linux-x86-64.so.2
7f94a6a0f000-7f94a6a11000 r--p 00033000 fd:01 45642507                   /usr/lib/x86_64-linux-gnu/ld-linux-x86-64.so.2
7f94a6a11000-7f94a6a13000 rw-p 00035000 fd:01 45642507                   /usr/lib/x86_64-linux-gnu/ld-linux-x86-64.so.2
7ffc422c4000-7ffc422e6000 rw-p 00000000 00:00 0                          [stack]
7ffc423a6000-7ffc423aa000 r--p 00000000 00:00 0                          [vvar]
7ffc423aa000-7ffc423ac000 r-xp 00000000 00:00 0                          [vdso]
ffffffffff600000-ffffffffff601000 --xp 00000000 00:00 0                  [vsyscall]`
	info, err := GetPythonProcInfo(bufio.NewScanner(bytes.NewReader([]byte(maps))))
	require.NoError(t, err)
	require.Equal(t, info.Musl, false)
	require.Equal(t, info.Version, PythonVersion{3, 6, 0})
	require.NotNil(t, info.PythonMaps)
	require.NotNil(t, info.LibPythonMaps)

	maps = `55c07863f000-55c078640000 r--p 00000000 00:32 39877587                   /usr/bin/python3.11
55c078640000-55c078641000 r-xp 00001000 00:32 39877587                   /usr/bin/python3.11
55c078641000-55c078642000 r--p 00002000 00:32 39877587                   /usr/bin/python3.11
55c078642000-55c078643000 r--p 00002000 00:32 39877587                   /usr/bin/python3.11
55c078643000-55c078644000 rw-p 00003000 00:32 39877587                   /usr/bin/python3.11
55c079447000-55c079448000 ---p 00000000 00:00 0                          [heap]
55c079448000-55c07944a000 rw-p 00000000 00:00 0                          [heap]
7f62329de000-7f6232a17000 rw-p 00000000 00:00 0 
7f6232a1e000-7f6232a26000 rw-p 00000000 00:00 0 
7f6232a27000-7f6232b2d000 rw-p 00000000 00:00 0 
7f6232b30000-7f6232b3e000 rw-p 00000000 00:00 0 
7f6232b3f000-7f6232b4a000 rw-p 00000000 00:00 0 
7f6232b4a000-7f6232b4b000 r--p 00000000 00:32 39878051                   /usr/lib/python3.11/lib-dynload/_opcode.cpython-311-x86_64-linux-musl.so
7f6232b4b000-7f6232b4c000 r-xp 00001000 00:32 39878051                   /usr/lib/python3.11/lib-dynload/_opcode.cpython-311-x86_64-linux-musl.so
7f6232b4c000-7f6232b4d000 r--p 00002000 00:32 39878051                   /usr/lib/python3.11/lib-dynload/_opcode.cpython-311-x86_64-linux-musl.so
7f6232b4d000-7f6232b4e000 r--p 00002000 00:32 39878051                   /usr/lib/python3.11/lib-dynload/_opcode.cpython-311-x86_64-linux-musl.so
7f6232b4e000-7f6232b4f000 rw-p 00003000 00:32 39878051                   /usr/lib/python3.11/lib-dynload/_opcode.cpython-311-x86_64-linux-musl.so
7f6232b4f000-7f6232b95000 rw-p 00000000 00:00 0 
7f6232b96000-7f6232c28000 rw-p 00000000 00:00 0 
7f6232c28000-7f6232c35000 r--p 00000000 00:32 39877572                   /usr/lib/libncursesw.so.6.4
7f6232c35000-7f6232c5f000 r-xp 0000d000 00:32 39877572                   /usr/lib/libncursesw.so.6.4
7f6232c5f000-7f6232c76000 r--p 00037000 00:32 39877572                   /usr/lib/libncursesw.so.6.4
7f6232c76000-7f6232c7b000 r--p 0004d000 00:32 39877572                   /usr/lib/libncursesw.so.6.4
7f6232c7b000-7f6232c7c000 rw-p 00052000 00:32 39877572                   /usr/lib/libncursesw.so.6.4
7f6232c7c000-7f6232c90000 r--p 00000000 00:32 39877577                   /usr/lib/libreadline.so.8.2
7f6232c90000-7f6232cb0000 r-xp 00014000 00:32 39877577                   /usr/lib/libreadline.so.8.2
7f6232cb0000-7f6232cb9000 r--p 00034000 00:32 39877577                   /usr/lib/libreadline.so.8.2
7f6232cb9000-7f6232cbc000 r--p 0003d000 00:32 39877577                   /usr/lib/libreadline.so.8.2
7f6232cbc000-7f6232cc2000 rw-p 00040000 00:32 39877577                   /usr/lib/libreadline.so.8.2
7f6232cc2000-7f6232cc3000 rw-p 00000000 00:00 0 
7f6232cc3000-7f6232cc5000 r--p 00000000 00:32 39878086                   /usr/lib/python3.11/lib-dynload/readline.cpython-311-x86_64-linux-musl.so
7f6232cc5000-7f6232cc7000 r-xp 00002000 00:32 39878086                   /usr/lib/python3.11/lib-dynload/readline.cpython-311-x86_64-linux-musl.so
7f6232cc7000-7f6232cc9000 r--p 00004000 00:32 39878086                   /usr/lib/python3.11/lib-dynload/readline.cpython-311-x86_64-linux-musl.so
7f6232cc9000-7f6232cca000 r--p 00006000 00:32 39878086                   /usr/lib/python3.11/lib-dynload/readline.cpython-311-x86_64-linux-musl.so
7f6232cca000-7f6232ccb000 rw-p 00007000 00:32 39878086                   /usr/lib/python3.11/lib-dynload/readline.cpython-311-x86_64-linux-musl.so
7f6232ccb000-7f6232dd4000 rw-p 00000000 00:00 0 
7f6232dd5000-7f6232faa000 rw-p 00000000 00:00 0 
7f6232faa000-7f6233016000 r--p 00000000 00:32 39877591                   /usr/lib/libpython3.11.so.1.0
7f6233016000-7f62331f8000 r-xp 0006c000 00:32 39877591                   /usr/lib/libpython3.11.so.1.0
7f62331f8000-7f62332e2000 r--p 0024e000 00:32 39877591                   /usr/lib/libpython3.11.so.1.0
7f62332e2000-7f6233311000 r--p 00338000 00:32 39877591                   /usr/lib/libpython3.11.so.1.0
7f6233311000-7f6233441000 rw-p 00367000 00:32 39877591                   /usr/lib/libpython3.11.so.1.0
7f6233441000-7f6233483000 rw-p 00000000 00:00 0 
7f6233483000-7f6233497000 r--p 00000000 00:32 39877100                   /lib/ld-musl-x86_64.so.1
7f6233497000-7f62334e3000 r-xp 00014000 00:32 39877100                   /lib/ld-musl-x86_64.so.1
7f62334e3000-7f6233519000 r--p 00060000 00:32 39877100                   /lib/ld-musl-x86_64.so.1
7f6233519000-7f623351a000 r--p 00095000 00:32 39877100                   /lib/ld-musl-x86_64.so.1
7f623351a000-7f623351b000 rw-p 00096000 00:32 39877100                   /lib/ld-musl-x86_64.so.1
7f623351b000-7f623351e000 rw-p 00000000 00:00 0 
7ffeabbad000-7ffeabbce000 rw-p 00000000 00:00 0                          [stack]
7ffeabbe1000-7ffeabbe5000 r--p 00000000 00:00 0                          [vvar]
7ffeabbe5000-7ffeabbe7000 r-xp 00000000 00:00 0                          [vdso]
ffffffffff600000-ffffffffff601000 --xp 00000000 00:00 0                  [vsyscall]`
	info, err = GetPythonProcInfo(bufio.NewScanner(bytes.NewReader([]byte(maps))))
	require.NoError(t, err)
	require.Equal(t, info.Musl, true)
	require.Equal(t, info.Version, PythonVersion{3, 11, 0})
	require.NotNil(t, info.PythonMaps)
	require.NotNil(t, info.LibPythonMaps)
}