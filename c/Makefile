CFLAGS = -g
LDFLAGS = -g -lreadline

SRCS = step0_repl.c
OBJS = $(SRCS:%.c=%.o)
BINS = $(OBJS:%.o=%)


build: $(word $(words $(BINS)),$(BINS))

$(OBJS): %.o: %.c
	gcc $(CFLAGS) -c $(@:%.o=%.c) -o $@

$(BINS): %: %.o
	gcc $+ -o $@ $(LDFLAGS)

clean:
	rm -f $(OBJS) $(BINS)
