NETWORK := scylladbtest

define get_scylla_ip
	$(shell docker inspect --format='{{ .NetworkSettings.Networks.$(NETWORK).IPAddress }}' $(1))
endef

run:
	go run main.go -hosts "$(call get_scylla_ip,scylla1),$(call get_scylla_ip,scylla2)"