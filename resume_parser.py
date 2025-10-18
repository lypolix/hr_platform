"""
NLP парсер резюме - извлечение навыков, технологий, опыта
АДАПТИРОВАН ДЛЯ ПРОМЫШЛЕННОГО СЕКТОРА (Технополис Москва)
"""
import re
import PyPDF2
import pdfplumber
from docx import Document
from pathlib import Path
from typing import Dict, List, Set, Tuple
import pandas as pd
from collections import Counter

# Импорт промышленных навыков
try:
    from industrial_skills import INDUSTRIAL_SKILLS, INDUSTRIAL_SPECIALIZATIONS
    USE_INDUSTRIAL = True
except ImportError:
    USE_INDUSTRIAL = False


class ResumeParser:
    """Парсер для извлечения данных из резюме (IT + Промышленность)"""
    
    # IT навыки (оригинальные)
    IT_SKILLS = {
        'languages': {
            'Python', 'Java', 'JavaScript', 'TypeScript', 'C++', 'C#', 'Go', 'Golang',
            'PHP', 'Ruby', 'Swift', 'Kotlin', 'Rust', 'Scala', 'R', 'SQL', 'HTML', 'CSS'
        },
        'frameworks_backend': {
            'Django', 'Flask', 'FastAPI', 'Spring', 'Spring Boot', 'Node.js', 'Express',
            'NestJS', 'Laravel', 'Symfony', 'Rails', 'ASP.NET', '.NET', 'Gin', 'Echo'
        },
        'frameworks_frontend': {
            'React', 'Vue', 'Angular', 'Svelte', 'Next.js', 'Nuxt', 'Redux', 'MobX',
            'jQuery', 'Bootstrap', 'Tailwind', 'Material-UI', 'Ant Design'
        },
        'databases': {
            'PostgreSQL', 'MySQL', 'MongoDB', 'Redis', 'Elasticsearch', 'SQLite',
            'Oracle', 'MS SQL', 'Cassandra', 'DynamoDB', 'ClickHouse', 'Kafka'
        },
        'devops': {
            'Docker', 'Kubernetes', 'Jenkins', 'GitLab CI', 'GitHub Actions', 'Terraform',
            'Ansible', 'AWS', 'Azure', 'GCP', 'CI/CD', 'Nginx', 'Apache'
        },
        'data_science': {
            'Pandas', 'NumPy', 'Scikit-learn', 'TensorFlow', 'PyTorch', 'Keras',
            'Jupyter', 'Matplotlib', 'Seaborn', 'OpenCV', 'NLTK', 'SpaCy', 'Transformers',
            'LLM', 'NLP', 'ML', 'Machine Learning', 'Deep Learning', 'Computer Vision'
        },
        'mobile': {
            'React Native', 'Flutter', 'iOS', 'Android', 'SwiftUI', 'Jetpack Compose'
        },
        'testing': {
            'Pytest', 'Jest', 'Selenium', 'Cypress', 'JUnit', 'TestNG', 'Postman'
        },
        'other': {
            'Git', 'REST API', 'GraphQL', 'gRPC', 'Microservices', 'Agile', 'Scrum',
            'Linux', 'Unix', 'WebSocket', 'OAuth', 'JWT', 'SOLID', 'OOP'
        }
    }
    
    # Объединяем IT и промышленные навыки
    TECH_SKILLS = IT_SKILLS.copy()
    if USE_INDUSTRIAL:
        TECH_SKILLS.update(INDUSTRIAL_SKILLS)
    
    # Категории специализаций (IT + Промышленность)
    IT_SPECIALIZATIONS = {
        'Backend': ['backend', 'бэкенд', 'бекенд', 'server', 'api'],
        'Frontend': ['frontend', 'фронтенд', 'react', 'vue', 'angular', 'web'],
        'Fullstack': ['fullstack', 'full-stack', 'full stack', 'фулстек'],
        'Data Science': ['data scientist', 'ml engineer', 'machine learning', 'nlp', 'аналитик'],
        'DevOps': ['devops', 'sre', 'infrastructure', 'kubernetes', 'docker'],
        'Mobile': ['mobile', 'ios', 'android', 'react native', 'flutter'],
        'QA': ['qa', 'quality assurance', 'тестировщик', 'test', 'manual qa', 'автотестирование'],
        'Product Manager': ['product manager', 'pm', 'продакт', 'менеджер продукта'],
        'Project Manager': ['project manager', 'руководитель проекта'],
    }
    
    # Объединяем IT и промышленные специализации
    SPECIALIZATIONS = IT_SPECIALIZATIONS.copy()
    if USE_INDUSTRIAL:
        SPECIALIZATIONS.update(INDUSTRIAL_SPECIALIZATIONS)
    
    def __init__(self):
        """Инициализация парсера"""
        self.all_skills = set()
        for category in self.TECH_SKILLS.values():
            self.all_skills.update([skill.lower() for skill in category])
    
    def extract_text_from_pdf(self, file_path: str) -> str:
        """Извлечение текста из PDF"""
        text = ""
        
        # Метод 1: PyPDF2
        try:
            with open(file_path, 'rb') as file:
                pdf_reader = PyPDF2.PdfReader(file)
                for page in pdf_reader.pages:
                    text += page.extract_text() or ""
        except Exception as e:
            print(f"PyPDF2 error for {file_path}: {e}")
        
        # Метод 2: pdfplumber (если PyPDF2 не справился)
        if len(text.strip()) < 100:
            try:
                with pdfplumber.open(file_path) as pdf:
                    for page in pdf.pages:
                        text += page.extract_text() or ""
            except Exception as e:
                print(f"pdfplumber error for {file_path}: {e}")
        
        return text
    
    def extract_text_from_docx(self, file_path: str) -> str:
        """Извлечение текста из DOCX"""
        try:
            doc = Document(file_path)
            return '\n'.join([para.text for para in doc.paragraphs])
        except Exception as e:
            print(f"DOCX error for {file_path}: {e}")
            return ""
    
    def extract_text(self, file_path: str) -> str:
        """Извлечение текста из файла"""
        file_path = Path(file_path)
        
        if file_path.suffix.lower() == '.pdf':
            return self.extract_text_from_pdf(str(file_path))
        elif file_path.suffix.lower() in ['.docx', '.doc']:
            return self.extract_text_from_docx(str(file_path))
        elif file_path.suffix.lower() == '.txt':
            # Чтение TXT с автоопределением кодировки
            try:
                # Пробуем UTF-8
                with open(file_path, 'r', encoding='utf-8') as f:
                    return f.read()
            except UnicodeDecodeError:
                # Пробуем Windows-1251 (кириллица)
                try:
                    with open(file_path, 'r', encoding='windows-1251') as f:
                        return f.read()
                except:
                    # Пробуем cp1252
                    try:
                        with open(file_path, 'r', encoding='cp1252') as f:
                            return f.read()
                    except:
                        return ""
        else:
            return ""
    
    def extract_skills(self, text: str) -> Dict[str, List[str]]:
        """Извлечение навыков из текста"""
        text_lower = text.lower()
        found_skills = {}
        
        for category, skills in self.TECH_SKILLS.items():
            category_skills = []
            for skill in skills:
                # Ищем навык с учетом границ слов
                pattern = r'\b' + re.escape(skill.lower()) + r'\b'
                if re.search(pattern, text_lower):
                    category_skills.append(skill)
            
            if category_skills:
                found_skills[category] = category_skills
        
        return found_skills
    
    def extract_experience_years(self, text: str) -> int:
        """Извлечение количества лет опыта"""
        patterns = [
            r'(\d+)\+?\s*(?:года|лет|years?)\s+(?:опыта|experience)',
            r'опыт[^\d]*(\d+)\s*(?:года|лет|years?)',
            r'experience[^\d]*(\d+)\s*years?',
        ]
        
        for pattern in patterns:
            match = re.search(pattern, text.lower())
            if match:
                return int(match.group(1))
        
        return 0
    
    def extract_salary(self, text: str) -> Tuple[int, int]:
        """Извлечение зарплаты из текста"""
        # Паттерны для поиска зарплаты
        patterns = [
            r'(?:зарплата|зп|оклад|salary)[:\s-]*(\d+)\s*(?:000)?[\s-]*(?:(\d+)\s*(?:000)?)?',
            r'(?:от|from)\s*(\d+)\s*(?:000|тыс|k)?\s*(?:до|to)?\s*(\d+)?\s*(?:000|тыс|k)?',
            r'(\d+)\s*(?:000|тыс|k)\s*[-–—]\s*(\d+)\s*(?:000|тыс|k)',
            r'(?:желаемая\s+зарплата|expected\s+salary)[:\s-]*(\d+)',
        ]
        
        for pattern in patterns:
            match = re.search(pattern, text.lower())
            if match:
                salary_from = int(match.group(1))
                # Если зарплата меньше 1000, значит указана в тысячах
                if salary_from < 1000:
                    salary_from *= 1000
                
                salary_to = None
                if match.lastindex >= 2 and match.group(2):
                    salary_to = int(match.group(2))
                    if salary_to < 1000:
                        salary_to *= 1000
                else:
                    salary_to = salary_from
                
                return (salary_from, salary_to)
        
        return (None, None)
    
    def detect_specialization(self, text: str, skills: Dict[str, List[str]]) -> List[str]:
        """Определение специализации"""
        text_lower = text.lower()
        detected = []
        
        for spec, keywords in self.SPECIALIZATIONS.items():
            for keyword in keywords:
                if keyword in text_lower:
                    detected.append(spec)
                    break
        
        # Дополнительная логика на основе навыков
        if not detected:
            if skills.get('frameworks_frontend') or skills.get('mobile'):
                if skills.get('frameworks_backend'):
                    detected.append('Fullstack')
                else:
                    detected.append('Frontend')
            elif skills.get('frameworks_backend') or skills.get('databases'):
                detected.append('Backend')
            elif skills.get('data_science'):
                detected.append('Data Science')
            elif skills.get('devops'):
                detected.append('DevOps')
        
        return list(set(detected))
    
    def parse_resume(self, file_path: str) -> Dict:
        """Полный парсинг резюме"""
        text = self.extract_text(file_path)
        
        if not text:
            return {
                'filename': Path(file_path).name,
                'text_length': 0,
                'skills': {},
                'all_skills': [],
                'specializations': [],
                'experience_years': 0,
                'salary_from': None,
                'salary_to': None,
                'error': 'Failed to extract text'
            }
        
        skills = self.extract_skills(text)
        all_skills_list = []
        for category_skills in skills.values():
            all_skills_list.extend(category_skills)
        
        specializations = self.detect_specialization(text, skills)
        experience = self.extract_experience_years(text)
        salary_from, salary_to = self.extract_salary(text)
        
        return {
            'filename': Path(file_path).name,
            'text_length': len(text),
            'skills': skills,
            'all_skills': all_skills_list,
            'specializations': specializations,
            'experience_years': experience,
            'salary_from': salary_from,
            'salary_to': salary_to,
            'text_preview': text[:500]
        }
    
    def parse_directory(self, directory_path: str) -> pd.DataFrame:
        """Парсинг всех резюме в директории"""
        directory = Path(directory_path)
        results = []
        
        for file_path in directory.iterdir():
            if file_path.suffix.lower() in ['.pdf', '.docx', '.doc', '.txt']:
                print(f"Parsing: {file_path.name}")
                result = self.parse_resume(str(file_path))
                results.append(result)
        
        return pd.DataFrame(results)
    
    def get_statistics(self, df: pd.DataFrame) -> Dict:
        """Получение статистики по всем резюме"""
        total_resumes = len(df)
        
        # Промышленные специализации для фильтрации
        industrial_specs = {
            'Логист', 'Токарь', 'Инженер-конструктор', 'Водитель погрузчика',
            'Слесарь', 'Мастер участка', 'Грузчик', 'Инженер по автоматизации',
            'Сварщик', 'Электромонтажник', 'Фрезеровщик', 'Инженер-технолог',
            'Кладовщик', 'Контролер ОТК', 'Инженер по качеству', 'Инженер по охране труда',
            'Инженер ПТО', 'Нормировщик', 'Слесарь-ремонтник', 'Наладчик оборудования'
        }
        
        # Подсчет только по промышленным специализациям
        spec_counter = Counter()
        for specs in df['specializations']:
            # Фильтруем только промышленные специализации
            industrial_only = [s for s in specs if s in industrial_specs]
            spec_counter.update(industrial_only)
        
        # Подсчет навыков
        skill_counter = Counter()
        for skills_list in df['all_skills']:
            skill_counter.update(skills_list)
        
        # Топ навыков
        top_skills = skill_counter.most_common(50)
        
        return {
            'total_resumes': total_resumes,
            'specializations': dict(spec_counter),
            'top_skills': top_skills,
            'avg_experience': df['experience_years'].mean(),
            'skills_distribution': skill_counter
        }


if __name__ == '__main__':
    # Тестирование
    parser = ResumeParser()
    resume_dir = r"c:\Users\Елена\Documents\Хакатоны\Моспром\Резюме"
    
    print("Начинаем парсинг резюме...")
    df = parser.parse_directory(resume_dir)
    
    print(f"\nОбработано резюме: {len(df)}")
    
    stats = parser.get_statistics(df)
    print(f"\nСтатистика:")
    print(f"Всего резюме: {stats['total_resumes']}")
    print(f"\nПо специализациям:")
    for spec, count in stats['specializations'].items():
        print(f"  {spec}: {count}")
    
    print(f"\nТоп-20 навыков:")
    for skill, count in stats['top_skills'][:20]:
        print(f"  {skill}: {count}")
    
    # Сохранение результатов
    df.to_json('parsed_resumes.json', orient='records', indent=2, force_ascii=False)
    print("\nРезультаты сохранены в parsed_resumes.json")
