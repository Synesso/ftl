package ftl.{{ .Name | camel | lower }}

import xyz.block.ftl.Context
import xyz.block.ftl.Ingress
import xyz.block.ftl.Method
import xyz.block.ftl.Verb

data class EchoRequest(val name: String? = "anonymous")
data class EchoResponse(val message: String)

class {{ .Name | camel }} {
  @Verb
  fun echo(context: Context, req: EchoRequest): EchoResponse {
    return EchoResponse(message = "Hello, ${req.name}!")
  }
}
