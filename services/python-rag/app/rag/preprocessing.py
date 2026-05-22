import re


class TextPreprocessor:
    """Normalizes raw text while preserving paragraph boundaries."""

    def clean(self, content: str) -> str:
        normalized = content.replace("\r\n", "\n").replace("\r", "\n")
        normalized = re.sub(r"[ \t]+", " ", normalized)
        normalized = re.sub(r"\n{3,}", "\n\n", normalized)
        return normalized.strip()
