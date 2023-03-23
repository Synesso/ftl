package xyz.block.ftl.drive

import org.eclipse.jetty.server.Server
import org.eclipse.jetty.server.ServerConnector
import org.eclipse.jetty.server.handler.ContextHandler
import org.eclipse.jetty.server.handler.HandlerList
import org.eclipse.jetty.server.handler.ResourceHandler
import org.eclipse.jetty.servlet.ServletHandler
import xyz.block.ftl.control.startControlChannel
import xyz.block.ftl.drive.transport.DriveServlet
import xyz.block.ftl.drive.verb.VerbDeck

val messages = listOf(
  "Warming up dilithium chamber...",
  "Initializing warp core...",
  "Sparking matter/anti-matter reactor...",
  "Engaging Proto-Star Drive...",
  "Connecting to the Mycelial Network..."
)

fun main(args: Array<String>) {
  Logging.init()
  val logger = Logging.logger("FTL Drive")
  logger.info(messages[(Math.random() * 10 % messages.size).toInt()])

  val server = Server()
  server.connectors = arrayOf(ServerConnector(server).apply {
    port = 8080
  })
  server.handler = HandlerList().apply {
    addHandler(ContextHandler("/_ftl").apply {
      handler = ResourceHandler().apply {
        resourceBase = "src/main/resources/web"
      }
    })
    addHandler(ServletHandler().apply {
      addServletWithMapping(DriveServlet::class.java, "/")
    })
  }

  VerbDeck.init("com.squareup.ftldemo")

  if (System.getenv("FTL_ENDPOINT") != null) {
    startControlChannel(logger, VerbDeck.instance)
  }
  server.start()
}
