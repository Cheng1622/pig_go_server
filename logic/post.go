package logic

import (
	"fmt"
	"go_server/dao/mysql"
	"go_server/dao/redis"
	"go_server/models"
	"go_server/pkg/snowflake"
	"strconv"

	"go.uber.org/zap"
)

func CreatePost(post *models.Post) (msg string, err error) {
	// 雪花算法 生成帖子id
	post.Id = snowflake.GenId()
	zap.L().Debug("createPostLogic", zap.Int64("postId", post.Id))
	err = mysql.InsertPost(post)
	if err != nil {
		return "failed", err
	}
	// 去点赞数量的 zset 新增一条记录
	err = redis.AddPost(post.Id)
	if err != nil {
		return "", err
	}
	//发表帖子成功时 要把帖子id 回给 请求方
	return strconv.FormatInt(post.Id, 10), nil
}

func GetPostList2(params *models.ParamListData) (apiPostDetailList []*models.ApiPostDetail, err error) {
	// 最热
	if params.Order == models.OrderByHot {
		// 先去redis 里面取 最新的数据
		ids, err := redis.GetPostIdsByScore(params.PageSize, params.PageNum)
		if err != nil {
			return nil, err
		}
		postLists, err := mysql.GetPostListByIds(ids)
		if err != nil {
			return nil, err
		}
		return rangeInitApiPostDetail(postLists)

	} else if params.Order == models.OrderByTime {
		//最新
		return GetPostList(params.PageSize, params.PageNum)
	}
	return nil, nil
}

// 分页全部
var offset int64

func GetPostList(pageSize int64, pageNum int64) (apiPostDetailList []*models.ApiPostDetail, err error) {

	offset = pageSize * (pageNum - 1)
	postList, err := mysql.GetPostList(offset, pageSize)
	if err != nil {
		return nil, err
	}
	return rangeInitApiPostDetail(postList)
}

// 分页按板块
func GetPostListByCommunity(community int64, pageSize int64, pageNum int64) (apiPostDetailList []*models.ApiPostDetail, err error) {

	offset = pageSize * (pageNum - 1)
	postList, err := mysql.GetPostListByCommunity(community, offset, pageSize)
	if err != nil {
		return nil, err
	}
	return rangeInitApiPostDetail(postList)
}

func rangeInitApiPostDetail(posts []*models.Post) (apiPostDetailList []*models.ApiPostDetail, err error) {
	for _, post := range posts {
		//再查 作者 名称
		email, err := mysql.GetEmailById(post.AuthorId)
		fmt.Println(post.Isnews)
		if post.Isnews == 0 && err != nil {
			zap.L().Warn("no author ")
			err = nil
			return nil, err
		}
		fmt.Println(post)
		//再查板块实体
		community, err := GetCommunityById(post.CommunityId)
		if err != nil {
			zap.L().Warn("no community ")
			err = nil
			return nil, err
		}
		apiPostDetail := new(models.ApiPostDetail)
		apiPostDetail.AuthorEmail = email
		apiPostDetail.Community = community
		apiPostDetail.Post = post
		apiPostDetailList = append(apiPostDetailList, apiPostDetail)
	}
	return apiPostDetailList, nil
}

func GetPostDetail(id int64) (apiPostDetail *models.ApiPostDetail, err error) {
	//先查帖子实体
	post, err := mysql.GetPostDetail(id)
	//再查 作者 名称
	email, err := mysql.GetEmailById(post.AuthorId)
	if post.Isnews == 0 && err != nil {
		zap.L().Warn("no author ")
		err = nil
	}
	//再查板块实体
	community, err := GetCommunityById(post.CommunityId)
	if err != nil {
		zap.L().Warn("no community ")
		err = nil
	}
	apiPostDetail = new(models.ApiPostDetail)
	apiPostDetail.AuthorEmail = email
	apiPostDetail.Community = community
	apiPostDetail.Post = post

	return apiPostDetail, err

}
