.PHONY:logic
logic:
	make logic -C logic

.PHONY:das
das:
	make das -C das

clean:
	make clean -C logic
	make clean -C das
