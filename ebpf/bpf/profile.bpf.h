#ifndef PROFILE_BPF_H
#define PROFILE_BPF_H


#define PERF_MAX_STACK_DEPTH      127
#define PROFILE_MAPS_SIZE         16384

#define KERN_STACKID_FLAGS (0 | BPF_F_FAST_STACK_CMP)
#define USER_STACKID_FLAGS (0 | BPF_F_FAST_STACK_CMP | BPF_F_USER_STACK)


struct sample_key {
    __u32 pid;
    __u32 flags;
    __s64 kern_stack;
    __s64 user_stack;
};

#define PROFILING_TYPE_UNKNOWN 1
#define PROFILING_TYPE_FRAMEPOINTERS 2
#define PROFILING_TYPE_PYTHON 3
#define PROFILING_TYPE_ERROR 4

struct pid_config {
    uint8_t type;
    uint8_t collect_user;
    uint8_t collect_kernel;
    uint8_t padding_;
};

#define OP_REQUEST_UNKNOWN_PROCESS_INFO 1

struct pid_event {
    uint32_t op;
    uint32_t pid;
};
struct pid_event e__;


struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __type(key, struct sample_key);
    __type(value, u32);
    __uint(max_entries, PROFILE_MAPS_SIZE);
} counts SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_STACK_TRACE);
    __uint(key_size, sizeof(u32));
    __uint(value_size, PERF_MAX_STACK_DEPTH * sizeof(u64));
    __uint(max_entries, PROFILE_MAPS_SIZE);
} stacks SEC(".maps");


struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __type(key, u32);
    __type(value, struct pid_config);
    __uint(max_entries, 1024);
} pids SEC(".maps");


struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(key_size, sizeof(u32));
    __uint(value_size, sizeof(u32));
} events SEC(".maps");




#endif // PROFILE_BPF_H