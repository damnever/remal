export

BINS = step0_repl step1_read_print step2_eval \
	   step3_env step4_if_fn_do step5_tco

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
