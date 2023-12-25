MUTOOL := mutool
XMLLINT := xmllint

GENERATED = $(patsubst %.pdf,%.xml,$(wildcard *.pdf))

all: $(GENERATED)

clean:
	rm -f $(GENERATED)

%.xml: %.pdf
	$(MUTOOL) draw -N -a -Fstext -- $< | $(XMLLINT) --format - > $@
