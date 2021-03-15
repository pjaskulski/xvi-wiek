./yaml2latex
mv xvi-wiek.tex ../tufte-latex/
cd ../tufte-latex/
pdflatex xvi-wiek.tex
pdflatex xvi-wiek.tex
texindy xvi-wiek.idx -L polish -C utf8
pdflatex xvi-wiek.tex
cp xvi-wiek.pdf ../xvi-wiek/ui/static/pdf/