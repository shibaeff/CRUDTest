deploy_region="westeurope"
deploy_group="go-custom-function"
az group create -l $deploy_region -n $deploy_group
function=$(az deployment group create --resource-group $deploy_group --template-uri https://raw.githubusercontent.com/groovy-sky/azure-func-go-handler/master/Template/azuredeploy.json | jq -r '. | .properties | .dependencies | .[] | .resourceName')
[ ! -d "azure-func-go-handler/.git" ] && git clone https://github.com/shibaeff/CRUDTest
cd CRUDTest/Function
go build *.go && func azure functionapp publish $function --no-build --force