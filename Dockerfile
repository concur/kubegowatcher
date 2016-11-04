FROM centurylink/ca-certs
# ADD ca-certificates.crt /etc/ssl/certs/
# build main w/ CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o main ./cnqrf5
ADD main /
CMD ["/main"]