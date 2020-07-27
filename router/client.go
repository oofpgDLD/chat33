package router

// Client is a middleman between the websocket connection and the hub
type Client struct {
	id         string
	device     string
	uuid       string
	createTime int64
	loginTime  int64
	user       *User

	conn Connection
}

func NewClient(id, device, uuid string, createTime, loginTime int64, conn Connection, u *User) *Client {
	return &Client{
		id:         id,
		device:     device,
		uuid:       uuid,
		conn:       conn,
		createTime: createTime,
		loginTime:  loginTime,
		user:       u,
	}
}

func (c *Client) Send(msg interface{}) {
	c.conn.WriteResponse(msg)
}

func (c *Client) Close(code int, msg string) error {
	return c.conn.Close(code, msg)
}

func (c *Client) GetId() string {
	return c.id
}

func (c *Client) GetDevice() string {
	return c.device
}

func (c *Client) GetUUID() string {
	return c.uuid
}

func (c *Client) GetCreateTime() int64 {
	return c.createTime
}

func (c *Client) GetLoginTime() int64 {
	return c.loginTime
}

func (c *Client) GetUser() *User {
	return c.user
}

func (c *Client) Run() error {
	return c.conn.Start()
}
