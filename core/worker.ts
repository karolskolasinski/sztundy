import { Hono } from "hono";
import postgres from "postgres";

interface Env {
  HYPERDRIVE: Hyperdrive;
}

const app = new Hono<{ Bindings: Env }>();

app.get("/", (c) => c.text("Hello Cloudflare!"));

app.get("/health/db", async (c) => {
  const sql = postgres(c.env.HYPERDRIVE.connectionString, {
    max: 5,
    fetch_types: false,
    prepare: true,
  });

  try {
    const [info] = await sql`
        SELECT current_database(), current_user
    `;

    return c.json({ connected: true, info });
  } catch (error) {
    console.error("PostgreSQL connection failed:", error);
    return c.json({ connected: false }, 500);
  }
});

export default app;