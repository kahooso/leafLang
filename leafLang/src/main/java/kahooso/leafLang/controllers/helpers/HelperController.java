package kahooso.leafLang.controllers.helpers;

import jakarta.transaction.Transactional;
import kahooso.leafLang.exceptions.NotFoundException;
import kahooso.leafLang.store.entities.WordEntity;
import kahooso.leafLang.store.repositories.WordRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Component;

@RequiredArgsConstructor
@Component
@Transactional
public class HelperController {

    private final WordRepository wordRepository;

    public WordEntity getWordOrThrowException(Long wordId) {

        return wordRepository
                .findById(wordId)
                .orElseThrow(() ->
                        new NotFoundException(
                                String.format(
                                        "Project with %s doesn't exists",
                                        wordId
                                )
                        )
                );
    }
}
