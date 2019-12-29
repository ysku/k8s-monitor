.PHONY: apply
apply:
	@kubectl apply -f serviceaccount.yml
	@kubectl apply -f clusterrolebinding.yml
	@kubectl apply -f deployment.yml


.PHONY: delete
delete:
	@kubectl delete -f serviceaccount.yml
	@kubectl delete -f clusterrolebinding.yml
	@kubectl delete -f deployment.yml

