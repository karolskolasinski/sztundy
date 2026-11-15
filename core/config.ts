export interface AppConfig {
  database: {
    host: string;
    port: number;
    user: string;
    password: string;
    database: string;
  };
  http: {
    port: number;
  };
}

function requiredEnv(name: string): string {
  const value = Deno.env.get(name);
  if (!value) {
    throw new Error(`Missing required environment variable: ${name}`);
  }

  return value;
}

export function loadConfig(): AppConfig {
  return {
    database: {
      host: requiredEnv("PG_HOST"),
      port: Number(requiredEnv("PG_PORT")),
      user: requiredEnv("PG_USER"),
      password: requiredEnv("PG_PASSWORD"),
      database: requiredEnv("PG_DATABASE"),
    },
    http: {
      port: Number(requiredEnv("CORE_PORT")),
    },
  };
}
