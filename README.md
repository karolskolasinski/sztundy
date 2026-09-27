# sztundy

## Czym jest SZTUNDY

SZTUNDY to aplikacja związana z organizacją treningów sportowych i pracą trenerów.

Początkowo system będzie przeznaczony przede wszystkim dla trenerów. W późniejszym etapie zostanie rozszerzony o pełną część zawodnika. System ma być ogólny względem dyscypliny sportowej. Pierwszym naturalnym zastosowaniem jest tenis stołowy, ale architektura nie powinna zakładać, że system jest przeznaczony wyłącznie dla tenisa stołowego. Docelowo system powinien nadawać się również np. dla:

- tenisa,
- padla,
- innych sportów, w których trener prowadzi indywidualne lub małe grupy treningowe i zarządza swoim grafikiem.

## Podstawowy problem biznesowy

Głównym problemem jest organizacja pracy trenera.

Trener ma wielu zawodników. Umawia z nimi treningi telefonicznie, osobiście, przez wiadomości lub w inny sposób. Zawodnicy mogą odwoływać treningi, zmieniać terminy albo pytać o dostępność. Trener nie przebywa przez cały dzień na obiekcie sportowym. Chce przyjeżdżać na konkretne bloki i prowadzić możliwie dużo treningów bez niepotrzebnych przerw. Problemem nie jest więc zwykłe prowadzenie kalendarza.

Kluczowy problem brzmi:

> Jak zaplanować treningi konkretnych zawodników w dostępnych blokach czasu trenera tak, aby grafik był możliwie zwarty i nie powodował niepotrzebnych dziur?

**Przykład:**

Trener jest dostępny:

```
06:00–10:00
15:00–18:00
```

Nie oznacza to jednej dostępności `06:00–18:00`. Są to dwa niezależne bloki dostępności.

Jeżeli w bloku `15:00–18:00` trener ma:

```
15:00–16:00 — trening,
17:00–18:00 — trening,
```

to godzina `16:00–17:00` jest potencjalnie użyteczną luką. Jeżeli system może zaplanować kolejny trening tak, aby wypełnić tę lukę, powinien preferować takie rozwiązanie. Z drugiej strony system nie powinien tworzyć małych, nieużytecznych luk, np. 30-minutowej przerwy, której nie da się wykorzystać na kolejny minimalny trening.

## Trener

Trener jest głównym użytkownikiem pierwszej wersji systemu.

- posiada listę swoich zawodników,
- może dodawać zawodników,
- może usuwać zawodników,
- może planować treningi,
- może zmieniać treningi,
- może anulować treningi,
- definiuje swoją dostępność,
- może mieć wiele niezależnych bloków dostępności w jednym dniu,
- może określać, które terminy udostępnia zawodnikom do samodzielnej rezerwacji.

Trener może prowadzić trening osobie posiadającej konto w SZTUNDY albo osobie, która konta nie posiada.

## Zawodnik

Zawodnik będzie drugą główną grupą użytkowników systemu, ale pełna część zawodnika zostanie zbudowana później. Nie wolno jednak projektować pierwszej wersji w sposób, który uniemożliwi późniejsze rozszerzenie systemu o zawodnika.

Docelowo zawodnik będzie mógł:

- posiadać konto,
- widzieć swoje treningi,
- mieć historię aktywności,
- widzieć historię treningów,
- analizować swoją aktywność,
- widzieć statystyki,
- śledzić progres,
- potencjalnie korzystać z dodatkowych funkcji związanych z treningiem,
- potencjalnie samodzielnie rezerwować treningi u trenera.

Ta część nie jest MVP pierwszej wersji, ale architektura systemu ma ją umożliwiać bez przebudowy istniejącego systemu.

## Ręcznie wpisany zawodnik a konto użytkownika

Trener może prowadzić trening osobie, która nie ma konta SZTUNDY. Może więc istnieć np. wpis: `„Karol Kowalski”.` Nie należy automatycznie próbować później powiązać historycznych wpisów z kontem użytkownika tylko dlatego, że nowy użytkownik ma takie samo imię i nazwisko. Jeżeli osoba później założy konto, historyczne ręczne wpisy pozostają historycznymi ręcznymi wpisami. Dopiero świadome wskazanie przez trenera konkretnego zarejestrowanego użytkownika powinno powodować powiązanie kolejnych treningów z tym użytkownikiem. System nie powinien zgadywać tożsamości.

## Dostępność trenera

Trener definiuje swoją dostępność jako zbiór niezależnych bloków.

Przykład:

```
Poniedziałek:
06:00–10:00
15:00–18:00

Wtorek:
08:00–12:00
16:00–20:00
```

Bloki są niezależnymi jednostkami. Nie należy traktować całego dnia jako jednego przedziału dostępności. Blok może mieć własną politykę dotyczącą pierwszego treningu.

Dostępne są trzy koncepcyjne tryby:

### Start from beginning

Pierwszy trening w pustym bloku powinien zaczynać się od początku bloku.

Przykład:

```
06:00–10:00
```

Preferowany pierwszy trening:

```
06:00–07:00
```

Ma to sens, gdy trener przyjeżdża specjalnie na trening i nie chce czekać.

### End at the end

Pierwszy trening może zostać dopasowany do końca bloku.

Przykład:

```
06:00–10:00
```

Dla treningu `08:00–10:00` taki termin może być preferowany.

### Flexible

Pierwszy trening może rozpocząć się w dowolnym poprawnym miejscu w obrębie bloku.

Przykład:

```
06:00–10:00
```

Dla treningu 2h możliwe są:

```
06:00–08:00
06:30–08:30
07:00–09:00
07:30–09:30
08:00–10:00
```

Polityka ta dotyczy przede wszystkim sytuacji, w której blok nie ma jeszcze zaplanowanych treningów.

## Czas treningu

W pierwszej wersji trening może trwać:

```
1h
1,5h
2h
2,5h
3h
```

Minimalny czas: `1h`, maksymalny czas: `3h`, krok: `30 minut`. To jest obecne założenie biznesowe i może zostać zmienione w przyszłości (na peweno musi być zmienione w przypadku np. obozów, zgrupowań).

## Planowanie treningu

Jednym z najważniejszych zachowań systemu jest: `Plan a Training`.  Trener chce zaplanować trening, więc może:

- wybrać zawodnika posiadającego konto
- wybrać zawodnika ręcznie wpisanego
- określić długość treningu
- otrzymać propozycje odpowiednich terminów
- wybrać termin
- zapisać trening

W przyszłości również zawodnik może rozpocząć ten proces (`Self-book a Training`). Zawodnik wybiera udostępnioną przez trenera dostępność i rezerwuje trening. Te dwa przypadki użycia mogą mieć wspólne mechanizmy biznesowe.

## Zarządzanie treningiem

Po zaplanowaniu treningu trener powinien móc:

- zmienić termin
- zmienić długość
- zmienić zawodnika
- anulować trening

Zmiana treningu może wpływać na dostępny grafik i musi respektować te same zasady dotyczące kolizji i zwartości grafiku. W przyszłości zawodnik może mieć możliwość odwołania lub zmiany treningu zgodnie z regułami systemu.

## Samodzielna rezerwacja zawodnika

Trener może zdecydować, że część jego dostępności będzie publicznie dostępna dla zawodników. Jeżeli trener wie jednak poza systemem, że konkretny zawodnik chce np. `środę 16:00`, system nie może tego wiedzieć. Trener powinien więc móc najpierw zaplanować lub zablokować ten termin. Termin zajęty przez zaplanowany trening nie powinien być jednocześnie oferowany jako wolna dostępność dla innych zawodników. System nie może rozwiązywać informacji, których nie posiada.

## Status treningu

Model powinien umożliwiać rozróżnienie różnych stanów cyklu życia treningu.

Przykładowo mogą istnieć stany odpowiadające sytuacjom:

- zaplanowany
- potwierdzony
- anulowany
- zakończony

## Przyszła część zawodnika

Po zbudowaniu części trenera system będzie rozszerzany o część zawodnika.

Docelowo może ona obejmować:

- historię treningów
- historię aktywności
- statystyki
- wykresy
- progres
- wyniki meczów
- turnieje
- aktywności sportowe
- inne dane treningowe

Ręczne wprowadzanie dużej ilości danych jest potencjalnie zbyt dużym obciążeniem dla użytkownika. Dlatego system powinien, gdzie to możliwe, pozyskiwać dane automatycznie.

## Automatyczne pozyskiwanie aktywności

Jednym z planowanych przyszłych źródeł danych jest ekosystem Google Health lub/i Apple. Docelowo chcemy wykorzystywać dane dotyczące aktywności użytkownika, np. aktywność typu:

- table tennis
- inne rozpoznawane aktywności sportowe

Interesują nas przede wszystkim:

- typ aktywności
- czas rozpoczęcia
- czas zakończenia
- czas trwania

Jeżeli użytkownik odbył trening zapisany wcześniej w SZTUNDY, system może w przyszłości zestawić informacje o zaplanowanym treningu z aktywnością pochodzącą z zewnętrznego źródła. Celem nie jest budowanie kolejnego ręcznego dzienniczka. Celem jest uzyskanie wartościowych danych przy minimalnym wysiłku użytkownika. Integracja z Google Health/Apple ma być traktowana jako zewnętrzne źródło/obszar volatility.

## Monetyzacja

Projekt może być rozwijany bezpłatnie, ale w przyszłości możliwe są:

- płatne funkcje,
- plan premium,
- inne modele monetyzacji.

Nie jest to jednak podstawowy problem, który mamy teraz rozwiązywać.

## Główne przypadki użycia:

1. Zarządzanie wydarzeniem sportowym (treningiem, sparingiem, obozem itp.)

Obejmuje:

- zarządzanie treningiem — umówienie/zmiana/odwołanie przez trenera/zawodnika
- zarządzanie sparingiem — umówienie/zmiana/odwołanie przez zawodnika
- zarządzanie dostępnością trenera
- zarządzanie dostępnością zawodnika
- wyszukiwanie wydarzenia sportowego
- dodawanie notatek do wydarzenia sportowego
- dodawanie oceny wydarzenia sportowego
- dodawanie/sugerowanie przez AI ćwiczeń do wydarzenia sportowego
- przeglądanie historii wydarzeń sportowych

2. Śledzenie i analiza statystyk, postępów oraz osiągnięć zawodnika

Obejmuje:

- przeglądanie statystyk — prezentacja danych, wykresy, wybór zakresów czasowych itd.
- pobieranie danych z różnych źródeł
- przetwarzanie/ujednolicanie danych
- wyliczanie statystyk, postępu, osiągnięć

## Pozostałe przypadki użycia:

1. zarządzanie profilem: informacje kontaktowe/sprzęt
2. zarządzanie użytkownikami (administrator aplikacji)
3. logowanie
4. płatności

## Obszary zmienności:

| #  | Obszar zmienności                           | Przykładowa zmienność                             |
|----|---------------------------------------------|---------------------------------------------------|
| 1  | **Typ i reguły wydarzenia**                 | trening, sparing, kolejne rodzaje aktywności      |
| 2  | **Workflow wydarzenia**                     | kto inicjuje, kolejność kroków, zmiana, odwołanie |
| 3  | **Dostępność uczestników**                  | dostępność trenera i zawodnika                    |
| 4  | **Reguły uczestnictwa**                     | indywidualne/grupowe, limity uczestników          |
| 5  | **Reguły planowania terminu**               | kolizje, luki, dopasowanie czasu, kompaktowość    |
| 6  | **Źródła danych aktywności**                | ręczne dane, urządzenia, zewnętrzne API           |
| 7  | **Transformacja danych**                    | mapowanie i ujednolicanie danych z różnych źródeł |
| 8  | **Reguły statystyk i postępu**              | sposób obliczania i agregowania danych            |
| 9  | **Klient / prezentacja**                    | trener, zawodnik, różne urządzenia i interfejsy   |
| 10 | **Powiadomienia**                           | odbiorcy, kanały, moment i sposób wysyłania       |
| 11 | **Security**                                | autentykacja, autoryzacja, credentials, sesje     |
| 12 | **Model i kontekst informacji dodatkowych** | notatki zawodnika, przeciwnika, aktywności itd.   |
