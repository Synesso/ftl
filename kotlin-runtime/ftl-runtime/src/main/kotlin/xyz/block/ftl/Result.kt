package xyz.block.ftl

import io.grpc.Status
import io.grpc.StatusException
import io.grpc.StatusRuntimeException

class Result<out T>(val value: Any?) {
  val isSuccess: Boolean get() = value !is Failure
  val isFailure: Boolean get() = value is Failure

  fun getOrNull(): T? =
    when {
      isFailure -> null
      else -> value as T
    }

  fun failureOrNull(): Status? =
    when (value) {
      is Failure -> value.status
      else -> null
    }

  override fun toString(): String =
    when (value) {
      is Failure -> value.toString() // "Failure($exception)"
      else -> "Success($value)"
    }

  data class Failure(val status: Status) {
    override fun equals(other: Any?): Boolean = other is Failure && status == other.status
    override fun hashCode(): Int = status.hashCode()
    override fun toString(): String = "Failure($status)"
  }

  companion object {
    /**
     * Returns an instance that encapsulates the given [value] as successful value.
     */
    fun <T> success(value: T): Result<T> = Result(value)

    /**
     * Returns an instance that encapsulates the given [Throwable] [exception] as failure.
     */
    fun failure(exception: Throwable): Result<Nothing> = Result(Failure(exception.toGrpcStatus()))

    private fun Throwable.toGrpcStatus(): Status {
      return when (this) {
        is StatusException -> this.status
        is StatusRuntimeException -> this.status
        else -> Status.INTERNAL.withDescription("Internal server error")
      }
    }
  }
}
