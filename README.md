# xvi-wiek.pl
Źródła powstającego powoli serwisu z bazą danych wydarzeń historycznych z XVI wieku. Serwis w wersji beta
powinien być widoczny pod tym adresem: [xvi-wiek.pl](http://xvi-wiek.pl).

Serwis wykorzystuje styl [Tufte CSS](https://edwardtufte.github.io/tufte-css/) (licencja MIT )z dodatkiem fontu [Crimson Pro](https://fonts.google.com/specimen/Crimson+Pro) (licencja Open Font License) Regular 400 i Regular 400 italic.

Backend serwisu obsługiwany jest przez program w języku Go, pliki yaml i mechanizm pamięci podręcznej typu klucz-wartość ([go-cache](https://github.com/patrickmn/go-cache), licencja MIT).
W roli routera wykorzystywany jest [chi](https://github.com/go-chi/chi).

Ikona Github pochodzi ze strony [GitHub Logos and Usage](https://github.com/logos).

Ikona Twitter pochodzi ze strony [Twitter Brand Resources](https://about.twitter.com/en_us/company/brand-resources.html).  

Dodatkowy mechanizm to bot wysyłający raz dziennie na konto [@xvi-wiek](https://twitter.com/xvi_wiek) informację o ciekawym wydarzeniu historycznym, które zdarzyło się w XVI wieku właśnie w tym dniu. Kod źródłowy bota dostępny jest w osobnym repozytorium: [xvi-wiek-bot](https://github.com/pjaskulski/xvi-wiek-bot). 

## API

    /api/dzien/{month}/{day}
  
    Zwraca wydarzenia z określonego dnia np. /api/dzien/4/10 zwróci
    wydarzenia z 10 kwietnia.

    /api/today
  
    Zwraca dzisiejsze wydarzenia, zapytanie uruchomione 1 stycznia 2020
    zwróci wydarzenia z 1 stycznia.   
