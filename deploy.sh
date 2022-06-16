#! /bin/bash
echo "###PARSING CONFIG FILES###"
sudo sed -i 's/localhost:2024/joangoma.dev/' views/layouts/cerberus.gohtml
sudo sed -i 's/localhost:2024/joangoma.dev/' config.yml
sudo sed -i 's/2024/443/' config/config.go

echo "###GO BUILD###"
sudo go build -buildvcs=false .
echo "###MOVING NECESSARY FILES###"
sudo mv jgt.solutions web
sudo mv web           ../webDeploy
sudo cp -r views      ../webDeploy
sudo cp -r assets     ../webDeploy
sudo cp -r css        ../webDeploy
sudo cp -r js         ../webDeploy
sudo cp -r vendor_web ../webDeploy
sudo cp -r fonts      ../webDeploy
sudo cp -r images     ../webDeploy
sudo cp config.yml    ../webDeploy
sudo cp cert.cer      ../webDeploy
sudo cp key.key       ../webDeploy
sudo cp runWeb.sh     ../webDeploy

echo "###TURNING GOOD CONFIG FOR LOCALHOST###"
sudo sed -i 's/joangoma.dev/localhost:2024/' views/layouts/cerberus.gohtml
sudo sed -i 's/joangoma.dev/localhost:2024/' config.yml
sudo sed -i 's/443/2024/' config/config.go

cd ..
ssh -tt -p24062 qathel@joangoma.dev  'sudo systemctl stop jgtweb.service'
scp -C -r -P 24062 webDeploy/* qathel@joangoma.dev:/srv/web
ssh -tt -p24062 qathel@joangoma.dev  'sudo systemctl start jgtweb.service'
