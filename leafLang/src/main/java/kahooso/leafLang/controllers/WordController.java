package kahooso.leafLang.controllers;

import jakarta.transaction.Transactional;
import kahooso.leafLang.Status;
import kahooso.leafLang.dto.AskDto;
import kahooso.leafLang.dto.WordDto;
import kahooso.leafLang.exceptions.BadRequestException;
import kahooso.leafLang.factories.WordDtoFactory;
import kahooso.leafLang.store.entities.WordEntity;
import kahooso.leafLang.store.repositories.WordRepository;
import kahooso.leafLang.controllers.helpers.HelperController;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Objects;
import java.util.Optional;
import java.util.stream.Collectors;
import java.util.stream.Stream;

@RestController
@Transactional
@RequiredArgsConstructor
public class WordController {

    private final WordRepository wordRepository;
    private final WordDtoFactory wordDtoFactory;

    public static final String GET_WORD = "api/words";
    public static final String DELETE_WORD = "api/words/{word_id}";
    public static final String CREATE_WORD = "api/words";
    public static final String UPDATE_WORD = "api/words/{word_id}";

    private final HelperController helperController;

    @GetMapping(GET_WORD)
    public List<WordDto> getWords(@RequestParam(value = "prefix_name", required = false) Optional<String> optionalPrefixName) {

        optionalPrefixName = optionalPrefixName.filter(prefixName -> !prefixName.trim().isEmpty());

        Stream<WordEntity> wordStream = optionalPrefixName
                .map(wordRepository::streamAllByInitialStartsWithIgnoreCase)
                .orElseGet(wordRepository::streamAllBy);

        return wordStream
                .map(wordDtoFactory::createWordDto)
                .collect(Collectors.toList());
    }

    @PostMapping(CREATE_WORD)
    public WordDto createWord(
            @RequestParam("word_initial") String wordInitial,
            @RequestParam("word_translate") String wordTranslate,
            @RequestParam("word_status") Status wordStatus,
            @RequestParam("word_example_of_usage") String wordExample_of_usage
    ) {

        if (wordInitial.trim().isEmpty()) {
            throw new BadRequestException("Word initial cannot be empty");
        }

        wordRepository.findByInitial(wordInitial)
                .ifPresent(word -> {
                    throw new BadRequestException(String.format("Project \"%s\" already exists.", wordInitial));
                });

        WordEntity word = wordRepository.saveAndFlush(
                WordEntity.builder()
                        .initial(wordInitial)
                        .translate(wordTranslate)
                        .example_of_usage(wordExample_of_usage)
                        .status(wordStatus)
                        .build()
        );

        return wordDtoFactory.createWordDto(word);
    }

    @PatchMapping(UPDATE_WORD)
    public WordDto updateWord(
            @PathVariable("word_id") long wordId,
            @RequestParam("word_initial") String wordInitial,
            @RequestParam("word_translate") String wordTranslate,
            @RequestParam("word_status") Status wordStatus,
            @RequestParam("word_example_of_usage") String wordExample_of_usage) {

        if (wordInitial.trim().isEmpty()) {
            throw new BadRequestException("Word initial cannot be empty");
        }

        WordEntity word = helperController.getWordOrThrowException(wordId);

        wordRepository
                .findByInitial(wordInitial)
                .filter(anotherProject -> !Objects.equals(anotherProject.getId(), wordId))
                .ifPresent(anotherProject -> {
                    throw new BadRequestException(String.format("Project \"%s\" already exists.", wordInitial));
                });

        word.setInitial(wordInitial);
        word.setTranslate(wordTranslate);
        word.setStatus(wordStatus);
        word.setExample_of_usage(wordExample_of_usage);

        wordRepository.saveAndFlush(word);

        return wordDtoFactory.createWordDto(word);
    }


    @DeleteMapping(DELETE_WORD)
    public AskDto deleteWord(@PathVariable("word_id") long wordId) {

        helperController.getWordOrThrowException(wordId);
        wordRepository.deleteById(wordId);
        return AskDto.makeDefault(true);
    }
}

/*
    HTTP спецификация

        Post - отвечает за создание чего-либо/запуск какой-либо логики
        Put - отвечает за полную замену/изменение объекта
        Patch - отвечает за обновление объекта
        Get - отвечает за получение какой-либо информации
        Delete - отвечает за удаление объекта
*/

