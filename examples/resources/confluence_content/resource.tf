resource "confluence_content" "content" {
  body  = "<h1>Contente created from terraform</h1><p>Hello</p>"
  space = "DEVOPS"
  title = "One interesting title"
  type  = "page"
}
