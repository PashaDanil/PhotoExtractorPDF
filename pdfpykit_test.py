from pikepdf import Pdf, Stream


def dump_pdf_complete(filename):
    """Полный дамп PDF структуры"""

    pdf = Pdf.open(filename)

    # Заголовок
    print(f'%PDF-{pdf.pdf_version}')
    print()

    # Все объекты
    for objnum in range(1, len(pdf.objects) + 1):
        try:
            obj = pdf.objects.get(objnum)
            if obj is None:
                continue

            print(f'{objnum} 0 obj')

            if isinstance(obj, Stream):
                # Это объект с потоком
                print(format_pikepdf_dict(obj))
                print('stream')

                try:
                    # Получаем распакованные данные
                    data = obj.read_bytes()
                    text = data.decode('latin-1', errors='replace')
                    print(text)
                except Exception as e:
                    print(f'[Ошибка чтения: {e}]')

                print('endstream')
            else:
                # Обычный объект
                print(format_pikepdf_obj(obj))

            print('endobj')
            print()

        except Exception as e:
            print(f'% Объект {objnum}: ошибка - {e}')
            print()

    # Trailer
    print('trailer')
    print(format_pikepdf_dict(pdf.trailer))
    print()
    print('%%EOF')

    pdf.close()


def format_pikepdf_obj(obj):
    """Форматирует объект pikepdf"""
    return str(obj)


def format_pikepdf_dict(d):
    """Форматирует словарь pikepdf"""
    lines = ['<<']
    for key, value in dict(d).items():
        lines.append(f'   {key} {value}')
    lines.append('>>')
    return '\n'.join(lines)


dump_pdf_complete('50137291M.pdf')