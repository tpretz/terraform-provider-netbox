test:
	go test -v $(shell go list ./... | grep -v /vendor/) 

testacc:
	TF_ACC=1 go test -v ./plugin/providers/netbox -run="TestAcc"

build: deps
	gox -osarch="linux/amd64 windows/amd64 darwin/amd64" \
	-output="pkg/{{.OS}}_{{.Arch}}/terraform-provider-netbox" .

install: clean build
	cp pkg/linux_amd64/terraform-provider-netbox ~/.terraform.d/plugins

tfplan: install
	terraform init -upgrade && TF_LOG=DEBUG terraform plan

tfapply: install
	terraform init -upgrade && TF_LOG=DEBUG terraform apply

release: release_bump release_build

release_bump:
	scripts/release_bump.sh

release_build:
	scripts/release_build.sh

deps:
#	go get -u github.com/hashicorp/terraform/plugin
	
clean:
	rm -rf pkg/
