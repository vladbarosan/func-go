azure-instance:
	./test/smoke.sh 1 1 2

local-instance:
	./test/smoke.sh 0 0 2

.PHONY: azure-instance local-instance
