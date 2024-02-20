# download kubebuilder and install locally.
curl -L -o kubebuilder "https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)"
sudo chmod +x kubebuilder && mv kubebuilder /usr/local/bin/

kubebuilder init --domain kcloutie --repo github.com/kcloutie/tekton-observer --owner kcloutie --project-name="tekton-observer" --component-config=true
kubebuilder create api --group config --version v1 --kind ControllerConfig --resource --controller=false --make=false

kubebuilder create api --group observer --version v1 --kind TektonObservation
