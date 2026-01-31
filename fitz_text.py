import fitz  # PyMuPDF


def dump_pdf_structure(filename):
    """Выводит структуру PDF используя PyMuPDF"""

    doc = fitz.open(filename)

    print(f'%PDF-{doc.metadata.get("format", "1.4")}')
    print()

    # Получаем все объекты
    xref_len = doc.xref_length()

    for i in range(1, xref_len):
        try:
            # Получаем сырое содержимое объекта
            obj_str = doc.xref_object(i)

            if obj_str:
                print(f'{i} 0 obj')
                print(obj_str)

                # Если есть stream, выводим его тоже
                try:
                    stream = doc.xref_stream(i)
                    if stream:
                        print('stream')
                        # Пробуем декодировать как текст
                        try:
                            print(stream.decode('latin-1'))
                        except:
                            print(f'[Бинарные данные: {len(stream)} байт]')
                        print('endstream')
                except:
                    pass

                print('endobj')
                print()
        except:
            pass

    doc.close()


dump_pdf_structure('50137291M.pdf')