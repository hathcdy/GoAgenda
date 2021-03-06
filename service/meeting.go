package service

import (
	"encoding/json"
	"github.com/chenf99/GoAgenda/entity"
	"log"
	"os"
)

//会议数据结构
type MeetingList struct {
	Meetings []entity.MeetingData `json:"meetings"`
}

var MeetingModel MeetingList
var meetingfile = os.Getenv("GOPATH") + "/src/github.com/chenf99/GoAgenda/data/meetings.json"

func init() {
	MeetingModel.readFromFile()
}

/**
* whether meeting exist
* @param title the title of meeting
* @return if success, true will be returned
 */
func (m *MeetingList) IsExist(title string) bool {
	for _, meeting := range m.Meetings {
		if meeting.Title == title {
			return true
		}
	}
	return false
}

/**
* get a meeting
* @param title the meeting's title
* @return if success, a meeting will be returned
 */
func (m *MeetingList) GetMeeting(title string) *entity.MeetingData {
	for _, meeting := range m.Meetings {
		if meeting.Title == title {
			return &meeting
		}
	}
	return &entity.MeetingData{}
}

/**
 * create a meeting
 * @param userName the sponsor's userName
 * @param title  the meeting's title
 * @param participator the meeting's participator
 * @param startData the meeting's start date
 * @param endData the meeting's end date
 * @return if success, true will be returned
 */
func (m *MeetingList) CreateMeeting(username, title, start, end string, participator []string) bool {
	meeting := entity.MeetingData{
		Title:        title,
		Sponsor:      username,
		Participator: participator,
		Start:        start,
		End:          end,
	}
	m.Meetings = append(m.Meetings, meeting)
	m.saveToFile()
	return true
}

/**
 * add a participator to a meeting
 * @param userName the sponsor's userName
 * @param title the meeting's title
 * @param participator the meeting's participator
 * @return if success, true will be returned
 */
func (m *MeetingList) AddMeetingParticipator(username, title, participator string) {
	for i, meeting := range m.Meetings {
		if meeting.Sponsor == username && meeting.Title == title && !meeting.IsParticipator(participator) {
			m.Meetings[i].AddParticipator(participator)
			m.saveToFile()
		}
	}
}

/**
 * remove a participator from a meeting
 * @param userName the sponsor's userName
 * @param title the meeting's title
 * @param participator the meeting's participator
 * @return if success, true will be returned
 */
func (m *MeetingList) RemoveMeetingParticipator(username, title, participator string) {
	for i, meeting := range m.Meetings {
		if meeting.Sponsor == username && meeting.Title == title && meeting.IsParticipator(participator) {
			m.Meetings[i].RemoveParticipator(participator)
			m.saveToFile()
		}
	}
}

/**
 * quit from a meeting
 * @param userName the current userName. no need to be the sponsor
 * @param title the meeting's title
 * @return if success, true will be returned
 */
func (m *MeetingList) QuitMeeting(username, title string) {
	for i, meeting := range m.Meetings {
		if meeting.Title == title && meeting.IsParticipator(username) {
			m.Meetings[i].RemoveParticipator(username)
			m.saveToFile()
		}
	}
}

/**
 * search a meeting by username, time interval
 * @param uesrName the sponsor's userName or as participator
 * @param startDate time interval's start date
 * @param endDate time interval's end date
 * @return a meeting list result
 */
func (m *MeetingList) MeetingQuery(username, start, end string) []entity.MeetingData {
	result := []entity.MeetingData{}
	for _, meeting := range m.Meetings {
		if meeting.Sponsor == username || meeting.IsParticipator(username) {
			if meeting.Start > start && meeting.End < end {
				result = append(result, meeting)
			}
		}
	}
	return result
}

/**
 * cancel a meeting by title and its sponsor
 * @param title meeting's title
 * @return if success, true will be returned
 */
func (m *MeetingList) CancelMeeting(title string) {
	for i, meeting := range m.Meetings {
		if meeting.Title == title {
			m.Meetings = append(m.Meetings[:i], m.Meetings[i+1:]...)
			m.saveToFile()
		}
	}
}

/**
 * delete all meetings by sponsor
 * @param userName sponsor's username
 * @return if success, true will be returned
 */
func (m *MeetingList) EmptyMeeting(userName string) {
	for _, meeting := range m.Meetings {
		if meeting.Sponsor == userName {
			m.CancelMeeting(meeting.Title)
		}
	}
	m.saveToFile()
}
func (m *MeetingList) saveToFile() {
	//MeetingList转json格式数据
	data, err := json.Marshal(*m)
	if err != nil {
		log.Fatal(err)
	}
	fp, err := os.OpenFile(meetingfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	//写入文件
	_, err = fp.Write(data)
	defer fp.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (u *MeetingList) readFromFile() {
	//判断文件是否存在
	_, err := os.Stat(meetingfile)
	if os.IsNotExist(err) {
		os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/chenf99/GoAgenda/data", 0777)
		return
	}
	fp, err := os.OpenFile(meetingfile, os.O_RDONLY, 0755)
	defer fp.Close()
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, 1000)
	//读取文件
	total, err := fp.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	//解析json数据到UserList
	err = json.Unmarshal(data[:total], u)
	if err != nil {
		log.Fatal(err)
	}
}
