resource "confluence_content" "content" {
  body  = "<h1>Contente created from terraform</h1>"
  space = "DEVOPS"
  title = "Example Title"
  type  = "page"
}

# Using template
resource "confluence_content" "content" {
  body  = templatefile("templates/page.tftpl", { name = "Example name" })
  space = "DEVOPS"
  title = "Example Title"
  type  = "page"
}

