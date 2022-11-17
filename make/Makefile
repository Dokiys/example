objects = foo.o bar.o
all: $(objects)
bar.o foo.o: %.o: %.m
	@echo $< && touch $@


# -----
include include_mk
.PHONY: b
b: a
	@echo $(B)


# -----
include var.mk
.PHONY: include_var
include_var:
	@echo "UNEXP_VAR: $(UNEXP_VAR)" && \
	echo "EXP_VAR: $(EXP_VAR)" && \
	echo "DEFEXP_VAR: $(DEFEXP_VAR)"


.PHONY: target_var
target_var: TARGET_VAR=1
target_var: TARGET_VAR2=2
target_var:
	@echo "TARGET_VAR: $(TARGET_VAR)\nTARGET_VAR2: $(TARGET_VAR2)"

.PHONY: t_target_var
t_target_var:
	@echo "TARGET_VAR: $(TARGET_VAR)"


# -----
BASH=/bin/bash
.PHONY: make_bash
make_bash:
	$(BASH) $(GFLAGS) '$(CMD)'

.PHONY: %
%:
	/bin/$* $(GFLAGS) '$(CMD)'

.PHONY: env_var
env_var:
	@echo "ENV_VAR: $(GOVERSION)"