package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func DB_Conect() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbSSLMode == "" {
		log.Fatal("Database environment variables not set")
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error conect from DB", err)
	}

	fmt.Println("DB connected")

	createTableProduct()
	createTableVideo()
	createTableBasket()
	createTablePurchaseRequest()
	createTablePurchaseItems()
	createTableSuccessfulPurchases()
	createPopulateFunction()

}

// сами курсы (4 шт)
func createTableProduct() {
	createTable := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		product_name TEXT NOT NULL,
		product_price DECIMAL(10,2) NOT NULL,
		currency TEXT DEFAULT 'RUB'
	);
	`
	_, err := DB.Exec(createTable)
	if err != nil {
		log.Fatal("Error created table product", err)
	}
	fmt.Println("Table products created successefully")
}

// 12 видео под каждый из курсов (4 курса, 48 видео)
func createTableVideo() {
	createTable := `
	CREATE TABLE IF NOT EXISTS video (
		id SERIAL PRIMARY KEY,
		product_id INTEGER NOT NULL,
		url TEXT NOT NULL,
		video_name TEXT NOT NULL
	);
	`
	_, err := DB.Exec(createTable)
	if err != nil {
		log.Fatal("Error created table video", err)
	}
	fmt.Println("Table video created successefully")
}

func createTablePurchaseRequest() {
	createTable := `
	CREATE TABLE IF NOT EXISTS purchase_request(
		id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL,
        email TEXT NOT NULL,
		total_amount DECIMAL(10,2) NOT NULL,
        payment_id TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := DB.Exec(createTable)
	if err != nil {
		log.Fatal("Error created table purchase_request", err)
	}
	fmt.Println("Table purchase_request created successefully")
}

func createTablePurchaseItems() {
	createTable := `
		CREATE TABLE IF NOT EXISTS purchase_item(
			id SERIAL PRIMARY KEY,
			purchase_request_id INTEGER NOT NULL REFERENCES purchase_request(id) ON DELETE CASCADE,
			product_id INTEGER NOT NULL,
			product_name TEXT NOT NULL,
			product_price DECIMAL(10,2) NOT NULL
		);
	`
	_, err := DB.Exec(createTable)
	if err != nil {
		log.Fatal("Error created table purchase_items", err)
	}
	fmt.Println("Table purchase_items created successefully")
}

func createTableSuccessfulPurchases() {
	createTable := `
		CREATE TABLE IF NOT EXISTS successful_purchases(
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			email TEXT NOT NULL,
			product_id INTEGER NOT NULL,
			product_name TEXT NOT NULL,
			product_price DECIMAL(10,2) NOT NULL,
			payment_id TEXT,
			sub_start TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			sub_end TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`
	_, err := DB.Exec(createTable)
	if err != nil {
		log.Fatal("Error created table successful_purchases", err)
	}
	fmt.Println("Table successful_purchases created successefully")
}

func createTableBasket() {
	createTable := `
		CREATE TABLE IF NOT EXISTS basket(
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			email TEXT NOT NULL,
			product_id INTEGER NOT NULL,
			product_name TEXT NOT NULL,
			product_price DECIMAL(10,2) NOT NULL
		);
	`
	_, err := DB.Exec(createTable)
	if err != nil {
		log.Fatal("Error created table basket", err)
	}
	fmt.Println("Table basket created successefully")
}

////////////////////////////////////////////////////////////////////////////

// createPopulateFunction создает SQL функцию для заполнения данных
func createPopulateFunction() error {
	sqlFunction := `
CREATE OR REPLACE FUNCTION populate_fitness_courses_func()
RETURNS VOID AS $$
DECLARE
    yoga_id INT;
    strength_id INT;
    cardio_id INT;
    pilates_id INT;
    video_counter INT;
    
    -- Названия видео для разных типов тренировок
    yoga_videos TEXT[] := ARRAY[
        'Введение в йогу', 'Позы для начинающих', 'Дыхательные практики',
        'Утренний комплекс', 'Йога для гибкости', 'Релаксация и медитация',
        'Йога для спины', 'Баланс и равновесие', 'Силовая йога',
        'Йога для сна', 'Продвинутые асаны', 'Итоговая практика'
    ];
    
    strength_videos TEXT[] := ARRAY[
        'Основы силового тренинга', 'Тренировка ног', 'Тренировка спины',
        'Грудные мышцы', 'Плечи и руки', 'Кор и пресс',
        'Функциональный тренинг', 'Тренировка с гантелями', 'Тренировка с резиной',
        'Прогрессия нагрузок', 'Восстановление', 'Итоговый комплекс'
    ];
    
    cardio_videos TEXT[] := ARRAY[
        'Кардио для начинающих', 'Интервальный тренинг', 'Танцевальное кардио',
        'HIIT тренировка', 'Скакалка', 'Лестница и степ',
        'Низкоударное кардио', 'Высокоинтенсивное кардио', 'Кардио для выносливости',
        'Тренировка на улице', 'Кардио с весом тела', 'Заминка и растяжка'
    ];
    
    pilates_videos TEXT[] := ARRAY[
        'Основы пилатеса', 'Пилатес для пресса', 'Пилатес для спины',
        'Работа с ковриком', 'Упражнения с роллом', 'Пилатес для ягодиц',
        'Пилатес для ног', 'Пилатес для осанки', 'Продвинутый уровень',
        'Пилатес для женщин', 'Пилатес для мужчин', 'Комплекс на все тело'
    ];
BEGIN
    -- Проверяем, есть ли уже курсы (очищаем если есть)
    DELETE FROM video;
    DELETE FROM products;
    
    -- Вставляем 4 фитнес-курса по одному с получением ID
    INSERT INTO products (product_name, product_price, currency)
    VALUES ('Йога для начинающих', 7999.99, 'RUB')
    RETURNING id INTO yoga_id;
    
    INSERT INTO products (product_name, product_price, currency)
    VALUES ('Силовые тренировки дома', 9999.50, 'RUB')
    RETURNING id INTO strength_id;
    
    INSERT INTO products (product_name, product_price, currency)
    VALUES ('Кардио тренировки', 6999.00, 'RUB')
    RETURNING id INTO cardio_id;
    
    INSERT INTO products (product_name, product_price, currency)
    VALUES ('Пилатес для осанки', 8999.75, 'RUB')
    RETURNING id INTO pilates_id;
    
    -- Добавляем по 12 видео для каждого курса
    FOR video_counter IN 1..12 LOOP
        -- Видео для йоги
        INSERT INTO video (product_id, url, video_name)
        VALUES (
            yoga_id,
            'https://fitness-academy.com/videos/yoga/video-' || video_counter,
            yoga_videos[video_counter]
        );
        
        -- Видео для силовых тренировок
        INSERT INTO video (product_id, url, video_name)
        VALUES (
            strength_id,
            'https://fitness-academy.com/videos/strength/video-' || video_counter,
            strength_videos[video_counter]
        );
        
        -- Видео для кардио
        INSERT INTO video (product_id, url, video_name)
        VALUES (
            cardio_id,
            'https://fitness-academy.com/videos/cardio/video-' || video_counter,
            cardio_videos[video_counter]
        );
        
        -- Видео для пилатеса
        INSERT INTO video (product_id, url, video_name)
        VALUES (
            pilates_id,
            'https://fitness-academy.com/videos/pilates/video-' || video_counter,
            pilates_videos[video_counter]
        );
    END LOOP;
    
    RAISE NOTICE '✅ Фитнес-курсы успешно созданы!';
    RAISE NOTICE '   Создано: 4 курса и 48 видео уроков';
END;
$$ LANGUAGE plpgsql;
`

	_, err := DB.Exec(sqlFunction)
	if err != nil {
		return fmt.Errorf("ошибка выполнения SQL: %w", err)
	}

	log.Println("✅ SQL функция populate_fitness_courses_func создана")
	return nil
}
