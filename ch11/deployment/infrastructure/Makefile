ECR_URL = $(shell terraform output -raw ecr_url)
REGION = $(shell terraform output -raw region)

ready: init validate

# Use this command to easily run all of the infrastructure deployment steps
# There are costs associated with this demo; Remember to tear everything down!
deploy: plan apply
	$(MAKE) kubeconfig
	@echo !!!
	@echo !!!
	@echo !!! There are cost associated with running this demo! Remember to tear everything down!
	@echo !!!
	@echo !!! Use 'make destroy' OR use 'terraform destroy' to tear everything down.
	@echo !!!
	@echo !!!

init:
	terraform init -upgrade

validate:
	terraform validate

plan:
	terraform plan -out infra.tfplan

apply:
	terraform apply infra.tfplan

destroy:
	terraform destroy -auto-approve

# Use this command to populate ~/.kube/config so that you can use k9s to browse the cluster
kubeconfig:
	aws eks update-kubeconfig --region ${REGION} --name $(shell terraform output project)
