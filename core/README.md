# Core

## Wymagania

- Go 1.23 lub nowszy
- PostgreSQL i istniejąca baza `sztundy_core`

## Konfiguracja

W katalogu głównym projektu wykonaj `cp .env.example .env`. W pliku `.env` ustaw właściwą nazwę roli i hasło PostgreSQL w `DATABASE_URL`. Jeśli Twoja rola nie nazywa się `sztundy`, zmień ją w tym adresie. Plik `.env` nie powinien trafiać do repozytorium.

Przykładowy adres dla lokalnego PostgreSQL:

`postgres://sztundy:HASLO_URL_ENCODED@127.0.0.1:5432/sztundy_core?sslmode=disable`

Jeśli hasło zawiera znaki specjalne używane w adresach URL, zakoduj je procentowo. Dla Neon użyj connection stringa z konsoli Neon i wymaganym `sslmode=require`.

Domyślnie serwer nasłuchuje na `127.0.0.1:8080`, aby nie udostępniać niedokończonego API w sieci. Adres można zmienić przez `HTTP_ADDR`.

## Schemat bazy danych

Baza `sztundy_core` musi już istnieć. Nie uruchamiaj ponownie `CREATE DATABASE`.

W pgAdmin wybierz bazę `sztundy_core`, otwórz Query Tool, wklej zawartość pliku `migrations/001_create_users.sql` i wykonaj zapytanie. Alternatywnie w terminalu z zainstalowanym `psql` uruchom:

`psql "postgres://sztundy:HASLO_URL_ENCODED@127.0.0.1:5432/sztundy_core?sslmode=disable" -f migrations/001_create_users.sql`

Tabela `users` przechowuje na tym etapie e-mail oraz pola profilu. Nie przechowuje hasła i nie implementuje logowania.

## Uruchomienie

Z katalogu głównego repozytorium wykonaj:

```bash
go mod tidy
go run ./cmd/core
```

Sprawdź serwer:

`curl http://127.0.0.1:8080/health`

Oczekiwana odpowiedź:

`{"status":"ok"}`

## API

Utworzenie profilu:

```bash
curl -i -X POST http://127.0.0.1:8080/api/profiles \
  -H 'Content-Type: application/json' \
  -d '{"email":"test@example.com","displayName":"Testowy gracz","bio":"Pierwszy profil"}'
```

Pobranie profilu; zamień `1` na `id` zwrócone podczas tworzenia:

`curl http://127.0.0.1:8080/api/profiles/1`

Aktualizacja wybranych pól:

```bash
curl -i -X PATCH http://127.0.0.1:8080/api/profiles/1 \
  -H 'Content-Type: application/json' \
  -d '{"displayName":"Nowa nazwa","bio":"Zaktualizowany opis"}'
```

## Odpowiedzialności warstw

`ProfileHandler` obsługuje HTTP, JSON i kody odpowiedzi. `ProfileManager` waliduje dane i realizuje przypadki użycia profilu. `UserAccess` wykonuje zapytania do PostgreSQL. `domain` zawiera model profilu i błędy domenowe. `database.Open` tworzy i sprawdza pulę połączeń. `migrations` zawiera schemat bazy.

To punkt startowy do lokalnego testowania przepływu. Przed wystawieniem API do Internetu trzeba dodać uwierzytelnianie i autoryzację: obecne endpointy nie sprawdzają, kto odczytuje ani zmienia profil. Domyślny adres `127.0.0.1` ogranicza dostęp do lokalnej maszyny.

## Zmienne środowiskowe

Konto `postgres` służy do administracji serwerem. Za jego pomocą możesz tworzyć bazy danych, role i nadawać uprawnienia. Natomiast aplikacja (**core**) powinna korzystać z osobnej roli, która ma uprawnienia do swojej bazy, ale nie ma pełnych uprawnień administratora serwera.

Po połączeniu przez pgAdmin otwórz Query Tool i wykonaj:
CREATE ROLE core WITH LOGIN PASSWORD 'TWOJE_MOCNE_HASLO';
CREATE DATABASE sztundy_db OWNER core;

Zastąp `TWOJE_MOCNE_HASLO` własnym, silnym hasłem. W rezultacie będziesz mieć konto administracyjne `postgres` oraz bazę `sztundy_db` z użytkownikiem aplikacyjnym `core`. To właśnie dane drugiego konta trafią później do `DATABASE_URL` w aplikacji.`