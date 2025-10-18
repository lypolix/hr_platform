"""
Flask API для NLP-анализа резюме
"""
from flask import Flask, request, jsonify, send_file, render_template
from flask_cors import CORS
import json
import os
from pathlib import Path
from resume_parser import ResumeParser
from resume_scorer import ResumeScorer
import pandas as pd
from datetime import datetime

app = Flask(__name__, static_folder='static', static_url_path='/static')
CORS(app)

# Конфигурация
RESUME_DIR = os.getenv('RESUME_DIR', str(Path(__file__).parent / "data" / "resumes_txt"))
DATA_DIR = Path(__file__).parent / "data"
DATA_DIR.mkdir(exist_ok=True)

# Глобальные объекты
parser = ResumeParser()
scorer = ResumeScorer()

# Кэш данных
cache = {
    'parsed_resumes': None,
    'statistics': None,
    'last_parse_time': None
}


@app.route('/')
def index():
    """Главная страница"""
    return render_template('index.html')


@app.route('/company')
def company():
    """Страница для представителя компании"""
    return render_template('company.html')


@app.route('/api/parse', methods=['POST'])
def parse_resumes():
    """Парсинг всех резюме в директории"""
    try:
        print("Начинаем парсинг резюме...")
        df = parser.parse_directory(RESUME_DIR)
        
        # Конвертируем в list of dicts
        resumes = df.to_dict('records')
        
        # Получаем статистику
        stats = parser.get_statistics(df)
        
        # Сохраняем в кэш
        cache['parsed_resumes'] = resumes
        cache['statistics'] = stats
        cache['last_parse_time'] = datetime.now().isoformat()
        
        # Сохраняем в файл
        output_file = DATA_DIR / 'parsed_resumes.json'
        with open(output_file, 'w', encoding='utf-8') as f:
            json.dump(resumes, f, ensure_ascii=False, indent=2)
        
        return jsonify({
            'success': True,
            'total_resumes': len(resumes),
            'statistics': stats,
            'message': f'Обработано {len(resumes)} резюме'
        })
    
    except Exception as e:
        return jsonify({
            'success': False,
            'error': str(e)
        }), 500


@app.route('/api/resumes', methods=['GET'])
def get_resumes():
    """Получение списка всех резюме"""
    try:
        # Фильтры
        specialization = request.args.get('specialization')
        skill = request.args.get('skill')
        min_experience = request.args.get('min_experience', type=int)
        
        resumes = cache['parsed_resumes']
        if not resumes:
            # Попытка загрузить из файла
            output_file = DATA_DIR / 'parsed_resumes.json'
            if output_file.exists():
                with open(output_file, 'r', encoding='utf-8') as f:
                    resumes = json.load(f)
                    cache['parsed_resumes'] = resumes
            else:
                return jsonify({
                    'success': False,
                    'error': 'Резюме не распарсены. Запустите /api/parse'
                }), 404
        
        # Применяем фильтры
        filtered = resumes
        
        if specialization:
            filtered = [r for r in filtered 
                       if specialization in r.get('specializations', [])]
        
        if skill:
            skill_lower = skill.lower()
            filtered = [r for r in filtered 
                       if skill_lower in [s.lower() for s in r.get('all_skills', [])]]
        
        if min_experience:
            filtered = [r for r in filtered 
                       if r.get('experience_years', 0) >= min_experience]
        
        # Добавляем информацию о вакансиях, на которые откликнулся кандидат
        applications_file = DATA_DIR / 'applications.json'
        vacancies_file = DATA_DIR / 'vacancies.json'
        
        if applications_file.exists() and vacancies_file.exists():
            with open(applications_file, 'r', encoding='utf-8') as f:
                applications = json.load(f)
            with open(vacancies_file, 'r', encoding='utf-8') as f:
                vacancies = json.load(f)
            
            # Создаем словарь вакансий для быстрого поиска
            vacancies_dict = {v['id']: v for v in vacancies}
            
            # Для каждого резюме находим вакансии, на которые оно откликалось
            for resume in filtered:
                resume_apps = [app for app in applications if app['resume_filename'] == resume.get('filename')]
                if resume_apps:
                    resume['applied_vacancies'] = []
                    for app in resume_apps:
                        vac = vacancies_dict.get(app['vacancy_id'])
                        if vac:
                            resume['applied_vacancies'].append({
                                'vacancy_id': vac['id'],
                                'title': vac['title'],
                                'company': vac['company'],
                                'match_score': app.get('match_score', 0)
                            })
        
        return jsonify({
            'success': True,
            'total': len(filtered),
            'resumes': filtered
        })
    
    except Exception as e:
        return jsonify({
            'success': False,
            'error': str(e)
        }), 500


@app.route('/api/statistics', methods=['GET'])
def get_statistics():
    """Получение статистики"""
    try:
        if cache['statistics']:
            return jsonify({
                'success': True,
                'statistics': cache['statistics'],
                'last_update': cache['last_parse_time']
            })
        
        # Попытка загрузить из файла
        output_file = DATA_DIR / 'parsed_resumes.json'
        if output_file.exists():
            with open(output_file, 'r', encoding='utf-8') as f:
                resumes = json.load(f)
            
            df = pd.DataFrame(resumes)
            stats = parser.get_statistics(df)
            cache['statistics'] = stats
            
            return jsonify({
                'success': True,
                'statistics': stats
            })
        
        return jsonify({
            'success': False,
            'error': 'Статистика не найдена. Запустите /api/parse'
        }), 404
    
    except Exception as e:
        return jsonify({
            'success': False,
            'error': str(e)
        }), 500


@app.route('/api/score', methods=['POST'])
def score_resumes():
    """Скоринг резюме по вакансии"""
    try:
        data = request.get_json()
        vacancy_text = data.get('vacancy_text')
        
        if not vacancy_text:
            return jsonify({
                'success': False,
                'error': 'Требуется текст вакансии (vacancy_text)'
            }), 400
        
        # Парсим требования вакансии
        vacancy_requirements = scorer.extract_requirements(vacancy_text)
        
        # Получаем резюме
        resumes = cache['parsed_resumes']
        if not resumes:
            output_file = DATA_DIR / 'parsed_resumes.json'
            if output_file.exists():
                with open(output_file, 'r', encoding='utf-8') as f:
                    resumes = json.load(f)
            else:
                return jsonify({
                    'success': False,
                    'error': 'Резюме не распарсены. Запустите /api/parse'
                }), 404
        
        # Ранжируем
        ranked = scorer.rank_resumes(resumes, vacancy_requirements)
        
        # Статистика по категориям
        category_stats = scorer.get_category_statistics(ranked)
        
        # Сохраняем результаты
        output_file = DATA_DIR / 'ranked_resumes.json'
        with open(output_file, 'w', encoding='utf-8') as f:
            json.dump(ranked, f, ensure_ascii=False, indent=2)
        
        return jsonify({
            'success': True,
            'total_resumes': len(ranked),
            'ranked_resumes': ranked,
            'category_statistics': category_stats,
            'vacancy_requirements': vacancy_requirements
        })
    
    except Exception as e:
        return jsonify({
            'success': False,
            'error': str(e)
        }), 500


@app.route('/api/export/json', methods=['GET'])
def export_json():
    """Экспорт результатов в JSON"""
    try:
        export_type = request.args.get('type', 'parsed')  # parsed или ranked
        
        if export_type == 'ranked':
            file_path = DATA_DIR / 'ranked_resumes.json'
        else:
            file_path = DATA_DIR / 'parsed_resumes.json'
        
        if not file_path.exists():
            return jsonify({
                'success': False,
                'error': f'Файл {export_type} не найден'
            }), 404
        
        return send_file(
            file_path,
            mimetype='application/json',
            as_attachment=True,
            download_name=f'{export_type}_resumes_{datetime.now().strftime("%Y%m%d_%H%M%S")}.json'
        )
    
    except Exception as e:
        return jsonify({
            'success': False,
            'error': str(e)
        }), 500


@app.route('/api/export/csv', methods=['GET'])
def export_csv():
    """Экспорт результатов в CSV"""
    try:
        export_type = request.args.get('type', 'parsed')
        
        if export_type == 'ranked':
            json_file = DATA_DIR / 'ranked_resumes.json'
        else:
            json_file = DATA_DIR / 'parsed_resumes.json'
        
        if not json_file.exists():
            return jsonify({
                'success': False,
                'error': f'Файл {export_type} не найден'
            }), 404
        
        # Загружаем JSON
        with open(json_file, 'r', encoding='utf-8') as f:
            data = json.load(f)
        
        # Создаем CSV
        df = pd.DataFrame(data)
        
        # Преобразуем списки в строки
        for col in df.columns:
            if df[col].dtype == 'object':
                df[col] = df[col].apply(lambda x: ', '.join(x) if isinstance(x, list) else str(x))
        
        csv_file = DATA_DIR / f'{export_type}_resumes.csv'
        df.to_csv(csv_file, index=False, encoding='utf-8-sig')
        
        return send_file(
            csv_file,
            mimetype='text/csv',
            as_attachment=True,
            download_name=f'{export_type}_resumes_{datetime.now().strftime("%Y%m%d_%H%M%S")}.csv'
        )
    
    except Exception as e:
        return jsonify({
            'success': False,
            'error': str(e)
        }), 500


@app.route('/api/skills/top', methods=['GET'])
def get_top_skills():
    """Получение топ навыков для облака тегов"""
    try:
        limit = request.args.get('limit', 50, type=int)
        
        # Проверяем кэш
        if cache['statistics'] and 'top_skills' in cache['statistics']:
            top_skills = cache['statistics']['top_skills'][:limit]
            return jsonify({
                'success': True,
                'skills': [{'name': skill, 'count': count} for skill, count in top_skills]
            })
        
        # Если кэша нет, генерируем из резюме
        parsed_file = DATA_DIR / 'parsed_resumes.json'
        if not parsed_file.exists():
            return jsonify({
                'success': False,
                'error': 'Данные резюме не найдены'
            }), 404
        
        with open(parsed_file, 'r', encoding='utf-8') as f:
            resumes = json.load(f)
        
        # Подсчитываем навыки
        skill_counts = {}
        for resume in resumes:
            for skill in resume.get('all_skills', []):
                skill_counts[skill] = skill_counts.get(skill, 0) + 1
        
        # Сортируем и берем топ
        top_skills = sorted(skill_counts.items(), key=lambda x: x[1], reverse=True)[:limit]
        
        return jsonify({
            'success': True,
            'skills': [{'name': skill, 'count': count} for skill, count in top_skills]
        })
    
    except Exception as e:
        return jsonify({
            'success': False,
            'error': str(e)
        }), 500


@app.route('/api/vacancies', methods=['GET'])
def get_vacancies():
    """Получение списка вакансий и статистики по зарплатам"""
    try:
        vacancies_file = DATA_DIR / 'vacancies.json'
        
        if not vacancies_file.exists():
            return jsonify({
                'success': False,
                'error': 'Файл с вакансиями не найден'
            }), 404
        
        with open(vacancies_file, 'r', encoding='utf-8') as f:
            vacancies = json.load(f)
        
        # Считаем статистику по зарплатам
        salary_by_spec = {}
        for vac in vacancies:
            spec = vac.get('specialization')
            salary = vac.get('salary_from')
            
            if spec and salary:
                if spec not in salary_by_spec:
                    salary_by_spec[spec] = []
                salary_by_spec[spec].append(salary)
        
        # Вычисляем средние
        stats = []
        for spec, salaries in salary_by_spec.items():
            avg = sum(salaries) / len(salaries)
            stats.append({
                'spec': spec,
                'count': len(salaries),
                'avg_salary': int(avg)
            })
        
        stats.sort(key=lambda x: -x['avg_salary'])
        
        return jsonify({
            'success': True,
            'vacancies': vacancies,
            'total': len(vacancies),
            'salary_stats': stats
        })
    
    except Exception as e:
        return jsonify({
            'success': False,
            'error': str(e)
        }), 500


@app.route('/api/companies', methods=['GET'])
def get_companies():
    """Получение статистики по компаниям"""
    try:
        # Загружаем вакансии
        vacancy_file = DATA_DIR / 'vacancies.json'
        if not vacancy_file.exists():
            return jsonify({
                'success': False,
                'error': 'Файл вакансий не найден'
            }), 404
        
        with open(vacancy_file, 'r', encoding='utf-8') as f:
            vacancies = json.load(f)
        
        # Считаем статистику по компаниям
        companies = {}
        for vac in vacancies:
            company = vac.get('company', 'Неизвестная компания')
            if company not in companies:
                companies[company] = {
                    'company': company,
                    'total_vacancies': 0,
                    'internships': 0,
                    'specializations': set()
                }
            
            companies[company]['total_vacancies'] += 1
            
            if vac.get('is_internship'):
                companies[company]['internships'] += 1
            
            spec = vac.get('specialization')
            if spec:
                companies[company]['specializations'].add(spec)
        
        # Преобразуем в список
        company_list = []
        for company_data in companies.values():
            company_list.append({
                'company': company_data['company'],
                'total_vacancies': company_data['total_vacancies'],
                'internships': company_data['internships'],
                'specializations': list(company_data['specializations']),
                'spec_count': len(company_data['specializations'])
            })
        
        # Сортируем по количеству вакансий
        company_list.sort(key=lambda x: x['total_vacancies'], reverse=True)
        
        # Общая статистика
        total_companies = len(company_list)
        total_vacancies = sum(c['total_vacancies'] for c in company_list)
        total_internships = sum(c['internships'] for c in company_list)
        
        return jsonify({
            'success': True,
            'companies': company_list,
            'summary': {
                'total_companies': total_companies,
                'total_vacancies': total_vacancies,
                'total_internships': total_internships
            }
        })
    
    except Exception as e:
        return jsonify({
            'success': False,
            'error': str(e)
        }), 500


@app.route('/api/resume/download/<filename>', methods=['GET'])
def download_resume(filename):
    """Скачивание резюме"""
    try:
        # Проверяем существование файла в data/resumes_txt
        resume_path = DATA_DIR / 'resumes_txt' / filename
        
        if not resume_path.exists():
            return jsonify({
                'success': False,
                'error': f'Файл {filename} не найден'
            }), 404
        
        return send_file(
            resume_path,
            as_attachment=True,
            download_name=filename,
            mimetype='text/plain'
        )
    
    except Exception as e:
        return jsonify({
            'success': False,
            'error': str(e)
        }), 500


@app.route('/api/applications', methods=['GET'])
def get_applications():
    """Получение откликов на вакансии"""
    try:
        applications_file = DATA_DIR / 'applications.json'
        
        if not applications_file.exists():
            return jsonify({
                'success': False,
                'error': 'Файл с откликами не найден'
            }), 404
        
        with open(applications_file, 'r', encoding='utf-8') as f:
            applications = json.load(f)
        
        # Фильтр по вакансии
        vacancy_id = request.args.get('vacancy_id')
        if vacancy_id:
            applications = [app for app in applications if app['vacancy_id'] == vacancy_id]
        
        return jsonify({
            'success': True,
            'applications': applications,
            'total': len(applications)
        })
    
    except Exception as e:
        return jsonify({
            'success': False,
            'error': str(e)
        }), 500


@app.route('/api/vacancies/with-applications', methods=['GET'])
def get_vacancies_with_applications():
    """Получение вакансий с количеством откликов, сгруппированных по отраслям"""
    try:
        vacancies_file = DATA_DIR / 'vacancies.json'
        applications_file = DATA_DIR / 'applications.json'
        
        if not vacancies_file.exists():
            return jsonify({
                'success': False,
                'error': 'Файл с вакансиями не найден'
            }), 404
        
        with open(vacancies_file, 'r', encoding='utf-8') as f:
            vacancies = json.load(f)
        
        # Загружаем отклики если есть
        applications = []
        if applications_file.exists():
            with open(applications_file, 'r', encoding='utf-8') as f:
                applications = json.load(f)
        
        # Считаем отклики по вакансиям
        applications_by_vacancy = {}
        for app in applications:
            vac_id = app['vacancy_id']
            applications_by_vacancy[vac_id] = applications_by_vacancy.get(vac_id, 0) + 1
        
        # Отрасли для группировки
        INDUSTRIES = {
            'Логистика и складское хозяйство': ['Логист', 'Водитель погрузчика', 'Грузчик'],
            'Машиностроение и металлообработка': ['Токарь', 'Фрезеровщик', 'Сварщик', 'Слесарь'],
            'Инженерно-технические специальности': ['Инженер-конструктор', 'Инженер-технолог', 'Инженер по автоматизации'],
            'Промышленная безопасность': ['Инженер по охране труда', 'Инженер по качеству'],
            'Электротехника': ['Электромонтажник'],
            'Производственное управление': ['Мастер участка', 'Кладовщик', 'Контролер ОТК']
        }
        
        # Группируем вакансии по отраслям
        vacancies_by_industry = {}
        total_vacancies = len(vacancies)
        total_applications = len(applications)
        
        for industry, specs in INDUSTRIES.items():
            industry_vacancies = []
            for vac in vacancies:
                if vac.get('specialization') in specs:
                    vac_with_apps = vac.copy()
                    vac_with_apps['applications_count'] = applications_by_vacancy.get(vac['id'], 0)
                    industry_vacancies.append(vac_with_apps)
            
            if industry_vacancies:
                vacancies_by_industry[industry] = {
                    'vacancies': industry_vacancies,
                    'count': len(industry_vacancies),
                    'total_applications': sum(v['applications_count'] for v in industry_vacancies)
                }
        
        return jsonify({
            'success': True,
            'total_vacancies': total_vacancies,
            'total_applications': total_applications,
            'industries': vacancies_by_industry
        })
    
    except Exception as e:
        return jsonify({
            'success': False,
            'error': str(e)
        }), 500


@app.route('/api/applications/<app_id>/status', methods=['PUT'])
def update_application_status(app_id):
    """Обновление статуса отклика"""
    try:
        data = request.get_json()
        new_status = data.get('status')
        
        # Валидация статуса
        valid_statuses = ['new', 'in_review', 'rejected', 'accepted']
        if new_status not in valid_statuses:
            return jsonify({
                'success': False,
                'error': f'Недопустимый статус. Разрешены: {", ".join(valid_statuses)}'
            }), 400
        
        applications_file = DATA_DIR / 'applications.json'
        
        if not applications_file.exists():
            return jsonify({
                'success': False,
                'error': 'Файл с откликами не найден'
            }), 404
        
        with open(applications_file, 'r', encoding='utf-8') as f:
            applications = json.load(f)
        
        # Находим отклик
        application = None
        for app in applications:
            if app.get('id') == app_id:
                application = app
                break
        
        if not application:
            return jsonify({
                'success': False,
                'error': f'Отклик с ID {app_id} не найден'
            }), 404
        
        # Обновляем статус
        old_status = application.get('status', 'new')
        application['status'] = new_status
        application['status_updated_at'] = datetime.now().isoformat()
        
        # Сохраняем обновленные данные
        with open(applications_file, 'w', encoding='utf-8') as f:
            json.dump(applications, f, ensure_ascii=False, indent=2)
        
        # Если статус изменен на "rejected", добавляем в общий банк резюме
        if new_status == 'rejected' and old_status != 'rejected':
            add_to_general_pool(application)
        
        return jsonify({
            'success': True,
            'application': application,
            'message': f'Статус изменен с "{old_status}" на "{new_status}"'
        })
    
    except Exception as e:
        return jsonify({
            'success': False,
            'error': str(e)
        }), 500


def add_to_general_pool(application):
    """Добавление отклоненного резюме в общий банк"""
    try:
        general_pool_file = DATA_DIR / 'general_resume_pool.json'
        parsed_resumes_file = DATA_DIR / 'parsed_resumes.json'
        
        # Загружаем общий банк
        if general_pool_file.exists():
            with open(general_pool_file, 'r', encoding='utf-8') as f:
                pool = json.load(f)
        else:
            pool = []
        
        # Проверяем, нет ли уже этого резюме в банке
        resume_filename = application.get('resume_filename')
        if any(r.get('filename') == resume_filename for r in pool):
            return  # Резюме уже в банке
        
        # Загружаем данные резюме из parsed_resumes.json
        if not parsed_resumes_file.exists():
            return
        
        with open(parsed_resumes_file, 'r', encoding='utf-8') as f:
            parsed_resumes = json.load(f)
        
        # Находим резюме
        resume_data = None
        for resume in parsed_resumes:
            if resume.get('filename') == resume_filename:
                resume_data = resume.copy()
                break
        
        if not resume_data:
            return
        
        # Добавляем метаданные
        resume_data['received_at'] = datetime.now().isoformat()
        resume_data['source'] = 'rejected_application'
        resume_data['status'] = 'available'
        resume_data['original_vacancy_id'] = application.get('vacancy_id')
        resume_data['rejection_date'] = application.get('status_updated_at')
        
        # Добавляем в банк
        pool.append(resume_data)
        
        # Сохраняем
        with open(general_pool_file, 'w', encoding='utf-8') as f:
            json.dump(pool, f, ensure_ascii=False, indent=2)
        
        print(f"Резюме {resume_filename} добавлено в общий банк (причина: отказ)")
    
    except Exception as e:
        print(f"Ошибка при добавлении в общий банк: {e}")


@app.route('/api/general-pool', methods=['GET'])
def get_general_pool():
    """Получение общего банка резюме"""
    try:
        general_pool_file = DATA_DIR / 'general_resume_pool.json'
        
        if not general_pool_file.exists():
            return jsonify({
                'success': True,
                'pool': [],
                'total': 0
            })
        
        with open(general_pool_file, 'r', encoding='utf-8') as f:
            pool = json.load(f)
        
        # Фильтры
        specialization = request.args.get('specialization')
        source = request.args.get('source')  # email, rejected_application
        
        filtered = pool
        
        if specialization:
            filtered = [r for r in filtered 
                       if specialization in r.get('specializations', [])]
        
        if source:
            filtered = [r for r in filtered if r.get('source') == source]
        
        return jsonify({
            'success': True,
            'pool': filtered,
            'total': len(filtered),
            'total_all': len(pool)
        })
    
    except Exception as e:
        return jsonify({
            'success': False,
            'error': str(e)
        }), 500


@app.route('/api/applications/stats', methods=['GET'])
def get_applications_stats():
    """Получение статистики по откликам (для компании или всех)"""
    try:
        company_name = request.args.get('company')
        
        applications_file = DATA_DIR / 'applications.json'
        vacancies_file = DATA_DIR / 'vacancies.json'
        
        if not applications_file.exists():
            return jsonify({
                'success': False,
                'error': 'Файл с откликами не найден'
            }), 404
        
        with open(applications_file, 'r', encoding='utf-8') as f:
            applications = json.load(f)
        
        # Если указана компания, фильтруем отклики по вакансиям компании
        if company_name:
            if not vacancies_file.exists():
                return jsonify({
                    'success': False,
                    'error': 'Файл с вакансиями не найден'
                }), 404
            
            with open(vacancies_file, 'r', encoding='utf-8') as f:
                vacancies = json.load(f)
            
            # Получаем ID вакансий компании
            company_vacancy_ids = [v['id'] for v in vacancies if v.get('company') == company_name]
            
            # Фильтруем отклики
            applications = [app for app in applications if app.get('vacancy_id') in company_vacancy_ids]
        
        # Считаем статистику по статусам
        stats = {
            'new': 0,
            'in_review': 0,
            'rejected': 0,
            'accepted': 0,
            'total': len(applications)
        }
        
        for app in applications:
            status = app.get('status', 'new')
            if status in stats:
                stats[status] += 1
        
        return jsonify({
            'success': True,
            'stats': stats,
            'company': company_name
        })
    
    except Exception as e:
        return jsonify({
            'success': False,
            'error': str(e)
        }), 500


@app.route('/api/applications/stats/all', methods=['GET'])
def get_all_companies_stats():
    """Получение статистики по всем компаниям"""
    try:
        applications_file = DATA_DIR / 'applications.json'
        vacancies_file = DATA_DIR / 'vacancies.json'
        
        if not applications_file.exists() or not vacancies_file.exists():
            return jsonify({
                'success': False,
                'error': 'Файлы данных не найдены'
            }), 404
        
        with open(applications_file, 'r', encoding='utf-8') as f:
            applications = json.load(f)
        
        with open(vacancies_file, 'r', encoding='utf-8') as f:
            vacancies = json.load(f)
        
        # Создаем словарь вакансий для быстрого поиска
        vacancies_dict = {v['id']: v for v in vacancies}
        
        # Группируем отклики по компаниям
        companies_stats = {}
        
        for app in applications:
            vacancy_id = app.get('vacancy_id')
            vacancy = vacancies_dict.get(vacancy_id)
            
            if not vacancy:
                continue
            
            company = vacancy.get('company', 'Неизвестная компания')
            status = app.get('status', 'new')
            
            if company not in companies_stats:
                companies_stats[company] = {
                    'company': company,
                    'new': 0,
                    'in_review': 0,
                    'rejected': 0,
                    'accepted': 0,
                    'total': 0
                }
            
            companies_stats[company]['total'] += 1
            if status in companies_stats[company]:
                companies_stats[company][status] += 1
        
        # Преобразуем в список
        stats_list = list(companies_stats.values())
        stats_list.sort(key=lambda x: x['total'], reverse=True)
        
        # Общая статистика
        total_stats = {
            'new': sum(c['new'] for c in stats_list),
            'in_review': sum(c['in_review'] for c in stats_list),
            'rejected': sum(c['rejected'] for c in stats_list),
            'accepted': sum(c['accepted'] for c in stats_list),
            'total': sum(c['total'] for c in stats_list)
        }
        
        return jsonify({
            'success': True,
            'companies': stats_list,
            'total_stats': total_stats
        })
    
    except Exception as e:
        return jsonify({
            'success': False,
            'error': str(e)
        }), 500


if __name__ == '__main__':
    import os
    
    # Получаем параметры из переменных окружения
    host = os.getenv('HOST', '0.0.0.0')
    port = int(os.getenv('PORT', 5000))
    debug = os.getenv('FLASK_ENV', 'development') == 'development'
    
    print("=" * 60)
    print("NLP Resume Analytics API")
    print("=" * 60)
    print(f"Директория с резюме: {RESUME_DIR}")
    print(f"Директория данных: {DATA_DIR}")
    print(f"Режим: {'Development' if debug else 'Production'}")
    print(f"Host: {host}")
    print(f"Port: {port}")
    print("\nEndpoints:")
    print("  POST /api/parse - Парсинг резюме")
    print("  GET  /api/resumes - Список резюме (фильтры: specialization, skill, min_experience)")
    print("  GET  /api/vacancies - Список вакансий и статистика по ЗП")
    print("  GET  /api/vacancies/with-applications - Вакансии с откликами")
    print("  GET  /api/applications - Список откликов (фильтр: vacancy_id)")
    print("  PUT  /api/applications/<id>/status - Обновление статуса отклика")
    print("  GET  /api/applications/stats - Статистика по откликам (фильтр: company)")
    print("  GET  /api/applications/stats/all - Статистика по всем компаниям")
    print("  GET  /api/general-pool - Общий банк резюме")
    print("  GET  /api/statistics - Статистика")
    print("  POST /api/score - Скоринг по вакансии")
    print("  GET  /api/export/json - Экспорт в JSON")
    print("  GET  /api/export/csv - Экспорт в CSV")
    print("  GET  /api/skills/top - Топ навыков")
    print("  GET  /api/resume/download/<filename> - Скачивание резюме")
    print("=" * 60)
    
    app.run(debug=debug, host=host, port=port)
