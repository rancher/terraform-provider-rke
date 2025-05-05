data "template_file" "addons" {
    template = file("${path.module}/manifest.yaml")

    vars = {
        do_token         = var.do_token
    }
}
