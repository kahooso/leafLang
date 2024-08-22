package kahooso.leafLang.factories;

import kahooso.leafLang.dto.WordDto;
import kahooso.leafLang.store.entities.WordEntity;
import org.springframework.stereotype.Component;

@Component
public class WordDtoFactory {

    public WordDto createWordDto(WordEntity wordEntity) {

        return WordDto
                .builder()
                .id(wordEntity.getId())
                .initial(wordEntity.getInitial())
                .translate(wordEntity.getTranslate())
                .status(wordEntity.getStatus())
                .example_of_usage(wordEntity.getExample_of_usage())
                .build();
    }
}