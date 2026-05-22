package study.rag.core;

import org.springframework.stereotype.Component;
import study.rag.core.model.SearchResult;
import study.rag.core.model.UserMemory;

import java.util.List;
import java.util.stream.Collectors;

@Component
public class ExtractiveAnswerGenerator {
  public String generate(String question, List<SearchResult> results, UserMemory memory) {
    if (results.isEmpty()) {
      return "I could not find matching knowledge base content yet. Add a document, index it, then ask again.";
    }
    String facts = memory.facts().isEmpty()
        ? ""
        : " I also considered your saved facts: " + String.join("; ", memory.facts()) + ".";
    String snippets = results.stream()
        .limit(3)
        .map(result -> result.chunk().content().length() > 420
            ? result.chunk().content().substring(0, 417).strip() + "..."
            : result.chunk().content())
        .collect(Collectors.joining("\n\n"));
    return "Based on the retrieved knowledge, here is the most relevant answer to: "
        + question
        + facts
        + "\n\n"
        + snippets;
  }
}
