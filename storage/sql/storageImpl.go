package sql

import (
	"github.com/jinzhu/gorm"
	"github.com/alexxxPopa/courses/conf"
	"github.com/pkg/errors"
	"github.com/alexxxPopa/courses/models"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Connection struct {
	db *gorm.DB
}

func (conn *Connection) Close() error {
	return conn.db.Close()
}

func (conn *Connection) Migrate() error {
	conn.db.DropTable(&models.User{}, &models.Plan{}, &models.Subscription{}) // for testing purposes
	conn.db = conn.db.AutoMigrate(&models.User{}, &models.Plan{}, &models.Subscription{})
	return conn.db.Error
}

func (conn *Connection) FindUserByEmail(email string) (*models.User, error) {
	return conn.findUser("email = ?", email)
}

func (conn *Connection) FindUserByStripeId(stripeId string) (*models.User, error) {
	return conn.findUser("stripe_id = ?", stripeId)
}

func (conn *Connection) CreateUser(user *models.User) error {
	tx := conn.db.Begin()
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (conn *Connection) findUser(query string, args ...interface{}) (*models.User, error) {
	user := &models.User{}
	values := append([]interface{}{query}, args...)
	if err := conn.db.First(user, values...).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (conn *Connection) UpdateUser(user *models.User) error {
	tx := conn.db.Begin()
	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (conn *Connection) FindPlanByTitle(title string) (*models.Plan, error) {
	plan := &models.Plan{}
	if err := conn.db.Model(&models.Plan{}).First(plan, "title = ?", title).Error; err != nil {
		return nil, err
	}
	return plan, nil
}

func (conn *Connection) FindPlans() ([]*models.PlanInfo, error) {
	plans := []*models.PlanInfo{}

	rows, _ := conn.db.Model(&models.Plan{}).
		Select("title, amount, currency, interval").Rows()
	defer rows.Close()
	for rows.Next() {
		plan := &models.PlanInfo{}
		err := rows.Scan(&plan.Title, &plan.Amount, &plan.Currency, &plan.Interval)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}

	return plans, nil
}

func (conn *Connection) CreatePlan(plan *models.Plan) error {
	tx := conn.db.Begin()
	if err := tx.Create(plan).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (conn *Connection) UpdatePlan(plan *models.Plan) error {
	tx := conn.db.Begin()
	if err := tx.Save(plan).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//func (conn *Connection) DeletePlan(plan *models.Plan) error {
//	err := conn.db.Delete(plan)
//	return err.Error
//}

func (conn *Connection) CreateSubscription(subscription *models.Subscription) error {
	tx := conn.db.Begin()
	if err := tx.Create(subscription).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (conn *Connection) FindSubscriptionByUser(user *models.User, status string) (*models.Subscription, error) {
	s := &models.Subscription{}
	rows, _ := conn.db.Model(&models.Subscription{}).Where("user_id =?", user.UserId).Rows()
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&s.SubscriptionId, &s.UserId, &s.PlanId, &s.StripeId, &s.Status, &s.Amount, &s.Currency, &s.PeriodStart, &s.PeriodEnd, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if s.Status == status {
			return s, nil
		}
	}
	return nil, nil
}

func (conn *Connection) UpdateSubscription(subscription *models.Subscription) error {
	tx := conn.db.Begin()
	if err := tx.Save(subscription).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (conn *Connection) IsSubscriptionActive(user *models.User, plan *models.Plan) bool {
	s := &models.Subscription{}
	err := conn.db.Model(&models.Subscription{}).Where(&models.Subscription{UserId:user.UserId, PlanId:plan.PlanId, Status:"Active"}).First(s).Error ; if err != nil {
		return false
	}
	return true
}

func (conn *Connection) CreateCourse(course *models.Course) error {
	tx := conn.db.Begin()
	if err := tx.Create(course).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (conn *Connection) GetCourses() ([]models.Course, error) {
	courses := &[]models.Course{}

	if err := conn.db.Find(courses).Error; err != nil {
		return nil, err
	}
	return *courses, nil
}

func (conn *Connection) CreateCategory(category *models.Category) error {
	tx := conn.db.Begin()
	if err := tx.Create(category).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (conn *Connection) FindCategoryById (categoryId uint) (*models.Category, error) {
	category := &models.Category{}
	if err := conn.db.Model(&models.Category{}).First(category, "category_id = ?", categoryId).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (conn *Connection) CreateArticle(article *models.Article) error {
	tx := conn.db.Begin()
	if err := tx.Create(article).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (conn *Connection) FindCourseById(courseId uint) (*models.Course, error) {
	course := &models.Course{}
	if err := conn.db.Model(&models.Course{}).First(course, "course_id = ?", courseId).Error; err != nil {
		return nil, err
	}
	return course, nil
}

func (conn *Connection) FindArticlesPerCourse(course *models.Course) ([]models.Article, error) {
	articles := &[] models.Article{}
	if err := conn.db.Model(course).Related(&models.Article{}, "course_id").Find(articles).Error ;  err!= nil {
		return nil, err
	}

	return *articles, nil
}

func Connect(config *conf.Config) (*Connection, error) {
	db, err := gorm.Open(config.DB.Driver, config.DB.Url)

	if err != nil {
		return nil, errors.Wrap(err, "Error while opening the database")
	}

	err = db.DB().Ping()

	if err != nil {
		return nil, errors.Wrap(err, "Error while connecting to the database")
	}

	conn := &Connection{
		db: db,
	}

	return conn, nil
}
