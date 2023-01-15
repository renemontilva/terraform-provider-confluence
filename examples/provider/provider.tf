# Using environment variables
# export CONFLUENCE_TOKEN=123token
# export CONFLUENCE_HOST=youruser.atlassian.net
# export CONFLUENCE_USER=youremail@example.com

provider "confluence" {

}

# Adding arguments to the provider block 
provider "confluence" {
  host  = "user.atlassian.net"
  user  = "user@example.com"
  token = "123token"

}
