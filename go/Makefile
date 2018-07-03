export

BINS = step0_repl step1_read_print

all: clean $(BINS)

mal: $(word $(words $(BINS)),$(BINS))
	cp $< $@

define build_template
$(1):
	env GOPATH=$(shell pwd) go build -o $$@ mal/cmd/$$@
endef

$(foreach b,$(BINS),$(eval $(call build_template,$(b))))

clean:
	rm -rf $(BINS)
