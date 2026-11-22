import postgres from "postgres";
import type { AppConfig } from "./config.ts";

export async function connectDatabase(config: AppConfig) {
  const isLocal = ["localhost", "127.0.0.1"].includes(config.database.host);

  const sql = postgres({
    host: config.database.host,
    port: config.database.port,
    user: config.database.user,
    password: config.database.password,
    database: config.database.database,
    // Neon requires TLS; the local PostgreSQL does not have it.
    ssl: isLocal ? false : "require",
    // The -pooler host is PgBouncer in transaction mode: disable prepared statements.
    prepare: false,
  });

  try {
    const [info] = await sql`SELECT current_database(), current_user, version(), now()`;
    console.info(`PostgreSQL connected at ${config.database.host}:${config.database.port}`, info);
    return sql;
  } catch (err) {
    console.error("PostgreSQL connection failed:", err);
    await sql.end({ timeout: 1 });
    throw err;
  }
}
