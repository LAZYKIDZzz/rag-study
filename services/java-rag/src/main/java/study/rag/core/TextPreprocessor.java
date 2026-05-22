package study.rag.core;

import org.springframework.stereotype.Component;

@Component
public class TextPreprocessor {
  public String clean(String content) {
    return content
        .replace("\r\n", "\n")
        .replace("\r", "\n")
        .replaceAll("[ \\t]+", " ")
        .replaceAll("\\n{3,}", "\n\n")
        .trim();
  }
}
