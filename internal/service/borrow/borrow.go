// ============================================================
// 我真诚地保证：
// 我自己独立地完成了整个程序从分析、设计到编码的所有工作。
// 如果在上述过程中，我遇到了什么困难而求教于人，那么，我将在程序实习报告中
// 详细地列举我所遇到的问题，以及别人给我的提示。
// 我的程序里中凡是引用到其他程序或文档之处，
// 例如教材、课堂笔记、网上的源代码以及其他参考书上的代码段,
// 我都已经在程序的注释里很清楚地注明了引用的出处。
// 我从未抄袭过别人的程序，也没有盗用别人的程序。
// 安俊豪
package borrow

import (
	"context"
	"errors"
	"library-management-api/internal/dao"
	"library-management-api/internal/model/entity"
	"library-management-api/internal/service/bizctx"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Service 借阅业务逻辑层
// Service 借阅业务结构体
type Service struct {
	bizCtxSvc *bizctx.Service
}

func New() *Service {
	return &Service{
		bizCtxSvc: bizctx.New(),
	}
}

// BorrowInput defines input for borrowing a book.
// BorrowInput 借书输入参数
type BorrowInput struct {
	UserId uint64
	BookId uint64
}

// BorrowOutput defines output of a borrow operation.
type BorrowOutput struct {
	Id       uint64
	BookName string
	BorrowAt string
}

// Borrow borrows a book for a user.
//借书（核心功能）
// 使用数据库事务保证数据一致性：
//   1. 检查图书是否存在且可借阅
//   2. 检查用户是否已借阅同一本书未还
//   3. 创建借阅记录
//   4. 更新图书状态为"已借出"
// 【重点】以上4步在同一个事务中，要么全部成功要么全部回滚
func (s *Service) Borrow(ctx context.Context, in BorrowInput) (*BorrowOutput, error) {
	var result BorrowOutput

		// 【数据库事务】保证借书操作的原子性
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 检查图书是否存在且可借阅
		var book entity.Books
		if err := tx.Model("books").Where("id", in.BookId).Scan(&book); err != nil {
			return err
		}
		if book.Id == 0 {
			return errors.New("图书不存在")
		}
				// 如果图书状态不是1（可借阅），则拒绝借出
		if book.Status != 1 {
			return errors.New("图书当前不可借阅")
		}
		result.BookName = book.Title

		// 检查用户是否有未归还的同本书
		count, err := tx.Model("borrows").
			Where("user_id", in.UserId).
			Where("book_id", in.BookId).
			Where("return_at IS NULL").
			Count()
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("您已借阅该书且尚未归还")
		}

		// 创建借阅记录
		id, err := tx.Model("borrows").Data(g.Map{
			"user_id": in.UserId,
			"book_id": in.BookId,
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		result.Id = uint64(id)

		// 更新图书状态为已借出
				// 步骤4：更新图书状态为"已借出(0)"
		if _, err := tx.Model("books").Where("id", in.BookId).Data(g.Map{"status": 0}).Update(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	customCtx := s.bizCtxSvc.Get(ctx)
	if customCtx != nil && customCtx.Session != nil {
		// 获取借阅时间（当前事务已提交）
		var borrow entity.Borrows
		if err := dao.Borrows.Ctx(ctx).Where("id", result.Id).Scan(&borrow); err == nil && borrow.BorrowAt != nil {
			result.BorrowAt = borrow.BorrowAt.String()
		}
	}

	return &result, nil
}

// Return marks a book as returned.
//还书（核心功能）
// 使用事务：
//   1. 检查借阅记录是否存在且未归还
//   2. 更新归还时间
//   3. 恢复图书状态为"可借阅"
func (s *Service) Return(ctx context.Context, borrowId uint64) (returnAt string, err error) {
		// 还书也使用事务保证数据一致性
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 查找借阅记录
		var borrow entity.Borrows
		if err := tx.Model("borrows").Where("id", borrowId).Scan(&borrow); err != nil {
			return err
		}
		if borrow.Id == 0 {
			return errors.New("借阅记录不存在")
		}
				// 如果已归还过，则拒绝重复操作
		if borrow.ReturnAt != nil {
			return errors.New("该书已归还")
		}

		// 更新归还时间
		_, err := tx.Model("borrows").
			Where("id", borrowId).
			Data(g.Map{"return_at": gtime.Now()}).Update()
		if err != nil {
			return err
		}

		// 恢复图书可借阅状态
				// 步骤3：恢复图书状态为"可借阅(1)"
		if _, err := tx.Model("books").Where("id", borrow.BookId).Data(g.Map{"status": 1}).Update(); err != nil {
			return err
		}

		returnAt = gtime.Now().String()
		return nil
	})

	return returnAt, err
}

// ListInput defines input for listing borrow records.
type ListInput struct {
	UserId uint64
	Page   int
	Size   int
	Status int // -1=全部, 0=借阅中, 1=已归还
}

// BorrowRecord represents a borrow record with user and book info.
type BorrowRecord struct {
	Id        uint64 `json:"id"`
	UserId    uint64 `json:"user_id"`
	UserName  string `json:"user_name"`
	BookId    uint64 `json:"book_id"`
	BookName  string `json:"book_name"`
	BorrowAt  string `json:"borrow_at"`
	ReturnAt  string `json:"return_at"`
}

//查询借阅列表（支持分页+状态筛选）
// 用 LEFT JOIN 关联 users 和 books 表，获取用户名和书名
func (s *Service) List(ctx context.Context, in ListInput) (total int, records []BorrowRecord, err error) {

    // 基础查询
    m := dao.Borrows.Ctx(ctx).
        LeftJoin("users", "borrows.user_id=users.id").
        LeftJoin("books", "borrows.book_id=books.id")

    // 用户筛选
    if in.UserId > 0 {
        m = m.Where("borrows.user_id", in.UserId)
    }

    // 状态筛选
    if in.Status == 0 {
        m = m.Where("borrows.return_at IS NULL")
    } else if in.Status == 1 {
        m = m.Where("borrows.return_at IS NOT NULL")
    }

    // 先统计总数（不要带 Fields）
    total, err = m.Count()
    if err != nil {
        return 0, nil, err
    }

    // 再查询数据
    err = m.Fields(
        "borrows.id",
        "borrows.user_id",
        "users.name as user_name",
        "borrows.book_id",
        "books.title as book_name",
        "borrows.borrow_at",
        "borrows.return_at",
    ).
        Page(in.Page, in.Size).
        OrderDesc("borrows.id").
        Scan(&records)

    if err != nil {
        return 0, nil, err
    }

    return total, records, nil
}

// GetCurrentUserId returns the currently logged-in user's ID.
func (s *Service) GetCurrentUserId(ctx context.Context) (uint64, error) {
	customCtx := s.bizCtxSvc.Get(ctx)
	if customCtx == nil || customCtx.User == nil {
		return 0, errors.New("用户未登录")
	}
	return customCtx.User.Id, nil
}



