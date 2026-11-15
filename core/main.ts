import { loadConfig } from "./config.ts";
import { connectDatabase } from "./database.ts";
import { Hono } from 'hono'

const config = loadConfig();
const db = await connectDatabase(config);
const app = new Hono()

app.get('/', (c) => c.text('Hello Deno!'))

Deno.serve({ port: config.http.port }, app.fetch)
