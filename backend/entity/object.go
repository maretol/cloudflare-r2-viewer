package entity

type Object struct {
	Path    string  // パス。/dir/file.ext という形で保存される
	Name    string  // ファイル名。 file.ext という形で保存される
	Size    int64   // サイズ。バイト数
	Sumnail *string // サムネイルのファイル名。nil時は未ダウンロード
}
