#Create an EC2 instance with Amazon Linux 2 with internet access
#Connect to your instance using putty
#Get the DB_url
echo DATABASE_URL=
read database_url

echo Server Running on PORT=
read port

#Perform a quick update on your instance:
sudo yum update -y

#Install Golang
sudo yum install -y golang

#Install git in your EC2 instance
sudo yum install git -y

sudo amazon-linux-extras install nginx1 -y

#Check git version
git version

#Cloning quotes-api / server from repo
git clone https://github.com/Skjaldbaka17/quotes-api.git
cd quotes-api
echo "DATABASE_URL=$database_url" > .env
echo "PORT=$port" >> .env
#start server https://www.tecmint.com/run-linux-command-process-in-background-detach-process/
nohup make server &
cd ..

#Need to create new nginx.conf in /etc/nginx looking like the nginx.conf in this root!
echo "Need to create new nginx.conf in /etc/nginx looking like the nginx.conf in this root!"

#Cloning db-setup from rep, do the setup from local comuter (if using free-tier or something with less than 4GB storage)
#git clone https://github.com/Skjaldbaka17/setup-quotes-db.git
#cd setup-quotes-db
#echo DATABASE_URL=
#read database_url
#echo "export DATABASE_URL=$database_url" > .env
#make setup
#cd ..



#if this fails maybe need to kill running http
#sudo fuser -k 80/tcp
#And then run start nginx again

# Get What process is using port 80
#fuser 80/tcp