apt-get update

apt-get install -y curl zip python3.7

curl -OL https://github.com/Azure/azure-functions-core-tools/releases/download/3.0.2106/Azure.Functions.Cli.linux-x64.3.0.2106.zip

unzip -d azure-functions-cli Azure.Functions.Cli.linux-x64.*.zip

rm Azure.Functions.Cli.linux-x64.*.zip

chmod +x azure-functions-cli/func

echo "alias func='/home/azure-functions-cli/func'" >> ~/.bashrc
