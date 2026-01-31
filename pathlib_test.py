from pathlib import Path
import zipfile
import fitz  # pip install pymupdf


def save_image_as_single_page_pdf(img_bytes: bytes, out_pdf: Path):
    """
    Создаёт одностраничный PDF, куда вставлена картинка целиком.
    Размер страницы подгоняется под размер картинки (в поинтах = пиксели).
    """
    img_doc = fitz.open("png", img_bytes)  # PyMuPDF понимает png/jpeg по сигнатуре
    pix = img_doc[0].get_pixmap(alpha=False)
    img_doc.close()

    doc = fitz.open()
    page = doc.new_page(width=pix.width, height=pix.height)
    page.insert_image(page.rect, stream=img_bytes)
    doc.save(str(out_pdf))
    doc.close()


def extract_embedded_images_only(pdf_path: str | Path) -> Path:
    pdf_file = Path(pdf_path)
    name = pdf_file.stem

    out_dir = pdf_file.parent / f"{name}_картинки"
    img_dir = out_dir / "images"
    pdf_dir = out_dir / "pdf"
    img_dir.mkdir(parents=True, exist_ok=True)
    pdf_dir.mkdir(parents=True, exist_ok=True)

    doc = fitz.open(str(pdf_file))

    for pno in range(len(doc)):
        page = doc.load_page(pno)

        # Только embedded images (Image XObject)
        images = page.get_images(full=True)  # список кортежей, где [0] = xref
        if not images:
            continue

        for i, img_info in enumerate(images, start=1):
            xref = img_info[0]
            extracted = doc.extract_image(xref)
            img_bytes = extracted["image"]
            ext = extracted.get("ext", "bin")

            # 1) Сохраняем как файл картинки
            out_img = img_dir / f"{name}_p{pno+1:02d}_i{i:02d}.{ext}"
            out_img.write_bytes(img_bytes)

            # 2) Сохраняем как PDF (одна картинка = один pdf-файл)
            out_pdf = pdf_dir / f"{name}_p{pno+1:02d}_i{i:02d}.pdf"
            save_image_as_single_page_pdf(img_bytes, out_pdf)

    doc.close()

    # ZIP папки
    zip_path = pdf_file.parent / f"{name}_картинки.zip"
    with zipfile.ZipFile(zip_path, "w", compression=zipfile.ZIP_DEFLATED) as z:
        for f in out_dir.rglob("*"):
            if f.is_file():
                z.write(f, f.relative_to(out_dir.parent))

    return zip_path


if __name__ == "__main__":
    # пример: python extract.py /path/to/50137291M.pdf
    import sys
    if len(sys.argv) < 2:
        print("Usage: python extract.py 50137291M.pdf")
        raise SystemExit(2)

    zip_path = extract_embedded_images_only(sys.argv[1])
    print(f"Done: {zip_path}")
