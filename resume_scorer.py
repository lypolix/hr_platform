"""
Скоринг резюме - сравнение с вакансией
"""
import re
from typing import Dict, List, Tuple
from collections import Counter
import numpy as np
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import cosine_similarity


class ResumeScorer:
    """Система скоринга резюме"""
    
    def __init__(self):
        """Инициализация"""
        self.vectorizer = TfidfVectorizer(
            lowercase=True,
            ngram_range=(1, 2),
            max_features=500
        )
    
    def extract_requirements(self, vacancy_text: str) -> Dict[str, List[str]]:
        """Извлечение требований из вакансии"""
        text_lower = vacancy_text.lower()
        
        # Секции с требованиями
        requirements_sections = [
            r'требования[:\s]+(.*?)(?=обязанности|условия|$)',
            r'requirements[:\s]+(.*?)(?=responsibilities|conditions|$)',
            r'необходимо[:\s]+(.*?)(?=обязанности|условия|$)',
            r'required[:\s]+(.*?)(?=responsibilities|conditions|$)',
        ]
        
        requirements_text = ""
        for pattern in requirements_sections:
            match = re.search(pattern, text_lower, re.DOTALL | re.IGNORECASE)
            if match:
                requirements_text += match.group(1)
        
        # Если не нашли секцию, используем весь текст
        if not requirements_text:
            requirements_text = text_lower
        
        # Извлекаем ключевые слова
        from resume_parser import ResumeParser
        parser = ResumeParser()
        skills = parser.extract_skills(requirements_text)
        
        # Извлекаем опыт
        experience = parser.extract_experience_years(requirements_text)
        
        return {
            'skills': skills,
            'all_skills': [skill for cat_skills in skills.values() for skill in cat_skills],
            'experience_required': experience,
            'text': requirements_text
        }
    
    def calculate_skills_match(self, resume_skills: List[str], required_skills: List[str]) -> Tuple[float, Dict]:
        """Расчет совпадения навыков"""
        resume_skills_lower = set([s.lower() for s in resume_skills])
        required_skills_lower = set([s.lower() for s in required_skills])
        
        if not required_skills_lower:
            return 0.0, {'matched': [], 'missing': []}
        
        matched = resume_skills_lower & required_skills_lower
        missing = required_skills_lower - resume_skills_lower
        
        match_score = len(matched) / len(required_skills_lower) * 100
        
        return match_score, {
            'matched': list(matched),
            'missing': list(missing),
            'matched_count': len(matched),
            'required_count': len(required_skills_lower)
        }
    
    def calculate_text_similarity(self, resume_text: str, vacancy_text: str) -> float:
        """Расчет текстового сходства через TF-IDF"""
        try:
            texts = [resume_text.lower(), vacancy_text.lower()]
            tfidf_matrix = self.vectorizer.fit_transform(texts)
            similarity = cosine_similarity(tfidf_matrix[0:1], tfidf_matrix[1:2])[0][0]
            return similarity * 100
        except:
            return 0.0
    
    def calculate_experience_match(self, resume_exp: int, required_exp: int) -> float:
        """Расчет совпадения опыта"""
        if required_exp == 0:
            return 100.0
        
        if resume_exp >= required_exp:
            return 100.0
        elif resume_exp >= required_exp * 0.7:
            return 80.0
        elif resume_exp >= required_exp * 0.5:
            return 60.0
        else:
            return 40.0
    
    def score_resume(self, resume_data: Dict, vacancy_data: Dict, weights: Dict = None) -> Dict:
        """
        Полный скоринг резюме
        
        Args:
            resume_data: данные резюме (из ResumeParser)
            vacancy_data: данные вакансии (из extract_requirements)
            weights: веса для компонентов скора
        
        Returns:
            Dict с результатами скоринга
        """
        if weights is None:
            weights = {
                'skills': 0.6,      # 60% - навыки
                'text': 0.2,        # 20% - текстовое сходство
                'experience': 0.2   # 20% - опыт
            }
        
        # Скор навыков
        skills_score, skills_details = self.calculate_skills_match(
            resume_data.get('all_skills', []),
            vacancy_data.get('all_skills', [])
        )
        
        # Текстовое сходство
        text_similarity = 0.0
        if 'text_preview' in resume_data and 'text' in vacancy_data:
            text_similarity = self.calculate_text_similarity(
                resume_data['text_preview'],
                vacancy_data['text']
            )
        
        # Скор опыта
        experience_score = self.calculate_experience_match(
            resume_data.get('experience_years', 0),
            vacancy_data.get('experience_required', 0)
        )
        
        # Итоговый скор
        total_score = (
            skills_score * weights['skills'] +
            text_similarity * weights['text'] +
            experience_score * weights['experience']
        )
        
        # Категоризация
        if total_score >= 80:
            category = 'excellent'
            color = 'green'
        elif total_score >= 60:
            category = 'good'
            color = 'yellow'
        elif total_score >= 40:
            category = 'average'
            color = 'orange'
        else:
            category = 'poor'
            color = 'red'
        
        return {
            'total_score': round(total_score, 1),
            'category': category,
            'color': color,
            'components': {
                'skills_score': round(skills_score, 1),
                'text_similarity': round(text_similarity, 1),
                'experience_score': round(experience_score, 1)
            },
            'skills_details': skills_details,
            'resume_filename': resume_data.get('filename', 'Unknown'),
            'explanation': self._generate_explanation(
                skills_details, 
                resume_data.get('experience_years', 0),
                vacancy_data.get('experience_required', 0)
            )
        }
    
    def _generate_explanation(self, skills_details: Dict, resume_exp: int, required_exp: int) -> str:
        """Генерация объяснения скора"""
        matched = skills_details.get('matched', [])
        missing = skills_details.get('missing', [])
        
        explanation_parts = []
        
        if matched:
            matched_str = ', '.join(list(matched)[:5])
            explanation_parts.append(f"Найдены навыки: {matched_str}")
            if len(matched) > 5:
                explanation_parts[-1] += f" и ещё {len(matched) - 5}"
        
        if missing:
            missing_str = ', '.join(list(missing)[:5])
            explanation_parts.append(f"Отсутствуют: {missing_str}")
            if len(missing) > 5:
                explanation_parts[-1] += f" и ещё {len(missing) - 5}"
        
        if required_exp > 0:
            if resume_exp >= required_exp:
                explanation_parts.append(f"Опыт {resume_exp} лет (требуется {required_exp})")
            else:
                explanation_parts.append(f"Опыт {resume_exp} лет (недостаточно, требуется {required_exp})")
        
        return '. '.join(explanation_parts) + '.' if explanation_parts else 'Недостаточно данных для анализа.'
    
    def rank_resumes(self, resumes_data: List[Dict], vacancy_data: Dict) -> List[Dict]:
        """Ранжирование всех резюме по релевантности"""
        scored_resumes = []
        
        for resume in resumes_data:
            score_result = self.score_resume(resume, vacancy_data)
            score_result['resume_data'] = resume
            scored_resumes.append(score_result)
        
        # Сортировка по убыванию скора
        scored_resumes.sort(key=lambda x: x['total_score'], reverse=True)
        
        return scored_resumes
    
    def get_category_statistics(self, scored_resumes: List[Dict]) -> Dict:
        """Статистика по категориям"""
        categories = Counter([r['category'] for r in scored_resumes])
        
        return {
            'total': len(scored_resumes),
            'excellent': categories.get('excellent', 0),
            'good': categories.get('good', 0),
            'average': categories.get('average', 0),
            'poor': categories.get('poor', 0),
            'categories_percent': {
                'excellent': round(categories.get('excellent', 0) / len(scored_resumes) * 100, 1),
                'good': round(categories.get('good', 0) / len(scored_resumes) * 100, 1),
                'average': round(categories.get('average', 0) / len(scored_resumes) * 100, 1),
                'poor': round(categories.get('poor', 0) / len(scored_resumes) * 100, 1),
            }
        }


if __name__ == '__main__':
    # Тестирование
    import json
    from resume_parser import ResumeParser
    
    # Пример вакансии
    vacancy_text = """
    Backend Python разработчик (Middle)
    
    Требования:
    - Опыт работы от 2 лет
    - Python, Django, Flask
    - PostgreSQL, Redis
    - Docker, Kubernetes
    - REST API, GraphQL
    - Git, Linux
    
    Будет плюсом:
    - Опыт с AWS
    - Знание микросервисной архитектуры
    """
    
    scorer = ResumeScorer()
    vacancy_requirements = scorer.extract_requirements(vacancy_text)
    
    print("Требования вакансии:")
    print(f"Навыки: {vacancy_requirements['all_skills']}")
    print(f"Опыт: {vacancy_requirements['experience_required']} лет")
    
    # Загружаем распарсенные резюме
    try:
        with open('parsed_resumes.json', 'r', encoding='utf-8') as f:
            resumes = json.load(f)
        
        print(f"\nРанжирование {len(resumes)} резюме...")
        ranked = scorer.rank_resumes(resumes, vacancy_requirements)
        
        print("\nТоп-10 кандидатов:")
        for i, result in enumerate(ranked[:10], 1):
            print(f"\n{i}. {result['resume_filename']}")
            print(f"   Скор: {result['total_score']}% ({result['category']})")
            print(f"   {result['explanation']}")
        
        # Статистика
        stats = scorer.get_category_statistics(ranked)
        print(f"\n\nСтатистика по категориям:")
        print(f"Отличные (80%+): {stats['excellent']} ({stats['categories_percent']['excellent']}%)")
        print(f"Хорошие (60-80%): {stats['good']} ({stats['categories_percent']['good']}%)")
        print(f"Средние (40-60%): {stats['average']} ({stats['categories_percent']['average']}%)")
        print(f"Слабые (<40%): {stats['poor']} ({stats['categories_percent']['poor']}%)")
        
        # Сохранение результатов
        with open('ranked_resumes.json', 'w', encoding='utf-8') as f:
            json.dump(ranked, f, ensure_ascii=False, indent=2)
        
        print("\nРезультаты сохранены в ranked_resumes.json")
        
    except FileNotFoundError:
        print("\nСначала запустите resume_parser.py для парсинга резюме!")
