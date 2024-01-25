package ftl.echo

import ftl.time.TimeModuleClient
import ftl.time.TimeRequest
import xyz.block.ftl.Context
import xyz.block.ftl.Ingress
import xyz.block.ftl.Method
import xyz.block.ftl.Verb
import xyz.block.ftl.Alias

class InvalidInput(val field: String) : Exception()

data class EchoRequest(@Alias("n") val name: String?)
data class EchoResponse(val message: String)

class Echo {
  @Throws(InvalidInput::class)
  @Verb
  @Ingress(Method.POST, "/echo")
  fun echo(context: Context, req: EchoRequest): EchoResponse {
    val response = context.call(TimeModuleClient::time, TimeRequest)
    return EchoResponse(message = "Hello, ${req.name ?: "anonymous"}! The time is ${response.time}.")
  }
}
