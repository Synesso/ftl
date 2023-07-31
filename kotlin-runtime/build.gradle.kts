buildscript {
  dependencies {
    classpath("com.squareup.wire:wire-gradle-plugin:4.7.2")
  }
}

group = "xyz.block"
version = "0.1.0-SNAPSHOT"

plugins {
  id("com.squareup.wire") version "4.7.2"
  kotlin("jvm") version "1.9.0"
  // Apply the java-library plugin for API and implementation separation.
  `java-library`
}

repositories {
  // Use Maven Central for resolving dependencies.
  mavenCentral()
}

dependencies {
  compileOnly(libs.hotswapAgentCore)

  // Use the Kotlin JUnit 5 integration.
  testImplementation(libs.kotlinTestJunit5)

  // Use the JUnit 5 integration.
  testImplementation(libs.junitJupiterEngine)
  testRuntimeOnly(libs.junitPlatformLauncher)

  // These dependencies are used internally, and not exposed to consumers on their own compile classpath.
  implementation(libs.classgraph)
  implementation(libs.logbackClassic)
  implementation(libs.logbackCore)
  implementation(libs.kotlinReflect)
  implementation(libs.kotlinxCoroutinesCore)
  implementation(libs.gson)
  implementation(libs.wireRuntime)
  implementation(libs.wireGrpcServer)
  implementation(libs.grpcNetty)
  implementation(libs.grpcProtobuf)
  implementation(libs.grpcStub)
}

// Disable gradlew because we use a Hermit-provided gradle.
tasks.findByName("wrapper")?.enabled = false

wire {
  kotlin {
    rpcRole = "server"
    rpcCallStyle = "blocking"
    grpcServerCompatible = true
  }
  sourcePath {
    srcDir("src/main/proto")
  }
}

tasks.named<Test>("test") {
  // Use JUnit Platform for unit tests.
  useJUnitPlatform()
  testLogging {
    events("passed", "skipped", "failed")
  }
}

tasks.jar {
  enabled = true
  isZip64 = true
  duplicatesStrategy = DuplicatesStrategy.EXCLUDE

  archiveFileName.set("${project.name}.jar")

  manifest {
    attributes["Main-Class"] = "xyz.block.ftl.main.MainKt"
  }

  from(sourceSets.main.get().output)
  dependsOn(configurations.compileClasspath)
  from({
    configurations.runtimeClasspath.get().filter { it.name.endsWith("jar") }.map { zipTree(it) }
  }) {
    exclude("META-INF/*.RSA", "META-INF/*.SF", "META-INF/*.DSA")
  }
}