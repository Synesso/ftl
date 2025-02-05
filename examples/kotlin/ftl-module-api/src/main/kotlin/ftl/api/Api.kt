package ftl.api

import com.google.gson.Gson
import com.google.gson.GsonBuilder
import com.google.gson.JsonDeserializer
import com.google.gson.JsonPrimitive
import com.google.gson.JsonSerializer
import ftl.builtin.HttpRequest
import ftl.builtin.HttpResponse
import xyz.block.ftl.*
import java.util.*
import java.util.concurrent.ConcurrentHashMap
import java.util.concurrent.atomic.AtomicInteger

data class Todo(
  val id: Int,
  val title: String,
  val completed: Boolean = false,
)

typealias GetStatusRequest = Unit

data class GetStatusResponse(
  val status: String,
)

data class GetTodoRequest(
  val id: Int,
)

data class GetTodoResponse(
  val todo: Todo?,
)

data class CreateTodoRequest(
  val title: String,
)

data class CreateTodoResponse(
  val id: Int,
)

fun makeGson(): Gson = GsonBuilder()
  .registerTypeAdapter(ByteArray::class.java, JsonSerializer<ByteArray> { src, _, _ ->
    JsonPrimitive(Base64.getEncoder().encodeToString(src))
  })
  .registerTypeAdapter(ByteArray::class.java, JsonDeserializer { json, _, _ ->
    Base64.getDecoder().decode(json.asString)
  })
  .create()

// These are global because Api is stateless.
private val todos = ConcurrentHashMap<Int, Todo>()
private val idCounter = AtomicInteger()
private val gson = makeGson()

class Api {
  private val headers = mapOf("Content-Type" to arrayListOf("application/json"))

  @Verb
  @Ingress(Method.GET, "/api/status")
  fun status(context: Context, req: HttpRequest<GetStatusRequest>): HttpResponse<GetStatusResponse> {
    return HttpResponse<GetStatusResponse>(status = 200, headers = mapOf(), body = GetStatusResponse("OK"))
  }

  @Verb
  @Ingress(Method.GET, "/api/todos/{id}")
  fun getTodo(context: Context, req: HttpRequest<GetTodoRequest>): HttpResponse<GetTodoResponse> {
    val todoId = req.pathParameters["id"]?.toIntOrNull()
    val todo = todos[todoId]

    return if (todo != null) {
      HttpResponse<GetTodoResponse>(
        status = 200,
        headers = mapOf(),
        body = GetTodoResponse(todo)
      )
    } else {
      HttpResponse<GetTodoResponse>(status = 404, headers = mapOf(), body = GetTodoResponse(null))
    }
  }

  @Verb
  @Ingress(Method.POST, "/api/todos")
  fun addTodo(context: Context, req: HttpRequest<CreateTodoRequest>): HttpResponse<CreateTodoResponse> {
    val todoReq = req.body
    val id = idCounter.incrementAndGet()
    todos.put(
      id, Todo(
        id = id,
        title = todoReq.title,
      )
    )

    return HttpResponse<CreateTodoResponse>(
      status = 201,
      headers = headers,
      body = CreateTodoResponse(id),
    )
  }
}
